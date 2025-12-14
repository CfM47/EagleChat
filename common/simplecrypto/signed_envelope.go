package simplecrypto

import "eaglechat/common/simplecrypto/rsa"

type SignedEnvelope struct {
	Content   []byte
	Signature []byte
}

func SignEnvelope(msg []byte, senderPrivKey *rsa.PrivateKey) (*SignedEnvelope, error) {
	signature, err := rsa.Sign(msg, senderPrivKey)
	if err != nil {
		return &SignedEnvelope{}, nil
	}

	return &SignedEnvelope{
		Content:   msg,
		Signature: signature,
	}, nil
}

func OpenSignedEnvelope(envelope *SignedEnvelope, senderPubKey *rsa.PublicKey) ([]byte, error) {
	err := rsa.Verify(envelope.Content, envelope.Signature, senderPubKey)
	if err != nil {
		return nil, err
	}

	return envelope.Content, nil
}
