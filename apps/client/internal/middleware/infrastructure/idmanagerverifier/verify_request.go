package idmanagerverifier

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	"eaglechat/common/ezlog"
	"eaglechat/common/simplecrypto/rsa"
)

type PublicKeyResponse struct {
	PublicKey []byte `json:"public_key"`
	Signature []byte `json:"signature"`
}

func VerifyIDManager(ctx context.Context, IP net.IP, port uint16, CAPubkey rsa.PublicKey) (*rsa.PublicKey, *http.Client, error) {
	baseURL := fmt.Sprintf("http://%s:%d", IP.String(), port)
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/pubkey", nil)
	if err != nil {
		msg := "Failed to create public key request"
		ezlog.Log(ctx).Errorf(msg+": %v", err)
		return nil, nil, fmt.Errorf(strings.ToLower(msg)+": %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		msg := "Failed to perform public key request"

		if urlErr, ok := err.(*url.Error); ok && urlErr.Timeout() {
			ezlog.Log(ctx).Warnf(msg+": %v", urlErr)
		} else {
			ezlog.Log(ctx).Errorf(strings.ToLower(msg)+": %v", urlErr)
		}

		return nil, nil, fmt.Errorf(strings.ToLower(msg)+": %w", err)
	}

	var response PublicKeyResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		msg := "Failed to decode public key request response"
		ezlog.Log(ctx).Errorf(msg+": v", err)
		return nil, nil, fmt.Errorf(strings.ToLower(msg)+": w", err)
	}

	err = rsa.Verify(response.PublicKey, response.Signature, &CAPubkey)
	if err != nil {
		msg := "ID Manager provided wrong signature for their public key"
		ezlog.Log(ctx).Warnf(msg+": %v", err)
		// no tolower to maintain ID as caps
		return nil, nil, fmt.Errorf(msg+": %w", err)
	}

	managerPK, err := rsa.PublicKeyFromBytes(response.PublicKey)
	if err != nil {
		msg := "ID Manager provided invalid signed public key"
		ezlog.Log(ctx).Errorf(msg+": %v", err)
		return nil, nil, fmt.Errorf(msg+": %w", err)
	}

	return managerPK, client, nil
}
