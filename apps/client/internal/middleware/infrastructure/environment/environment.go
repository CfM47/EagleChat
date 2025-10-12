package environment

import (
	middleware_entities "eaglechat/apps/client/internal/middleware/domain/entities"
	"eaglechat/common/simplecrypto/rsa"
	"fmt"
	"net"
	"os"
	"strconv"
)

func GetDefaultIDManagerData() (middleware_entities.IDManagerData, error) {
	ipEnv, ok := os.LookupEnv("ID_MANAGER_IP")
	if !ok {
		return middleware_entities.IDManagerData{}, fmt.Errorf("no default ID manager configured: missing ID_MANAGER_IP")
	}
	ip := net.ParseIP(ipEnv)
	if ip == nil {
		return middleware_entities.IDManagerData{}, fmt.Errorf("invalid ID_MANAGER_IP: %s", ipEnv)
	}

	portEnv, ok := os.LookupEnv("ID_MANAGER_PORT")
	if !ok {
		return middleware_entities.IDManagerData{}, fmt.Errorf("no default ID manager configured: missing ID_MANAGER_PORT")
	}
	port, err := strconv.ParseUint(portEnv, 10, 16)
	if err != nil {
		return middleware_entities.IDManagerData{}, fmt.Errorf("invalid ID_MANAGER_PORT: %s", portEnv)
	}

	pkFile, ok := os.LookupEnv("ID_MANAGER_PUBLIC_KEY_FILE")
	if !ok {
		return middleware_entities.IDManagerData{}, fmt.Errorf("no default ID manager configured: missing ID_MANAGER_PUBLIC_KEY_FILE")
	}

	pkBytes, err := os.ReadFile(pkFile)
	if err != nil {
		return middleware_entities.IDManagerData{}, fmt.Errorf("failed to read public key file: %w", err)
	}
	pk, err := rsa.PublicKeyFromBytes(pkBytes)
	if err != nil {
		return middleware_entities.IDManagerData{}, fmt.Errorf("invalid public key in ID_MANAGER_PUBLIC_KEY_FILE: %w", err)
	}

	return middleware_entities.NewIDManagerData(ip, uint16(port), *pk), nil
}
