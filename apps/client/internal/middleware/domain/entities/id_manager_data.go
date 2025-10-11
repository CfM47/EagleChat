package entities

import (
	"eaglechat/apps/client/internal/utils/simplecrypto/rsa"
	"net"
)

type IDManagerData struct {
	IP        net.IP
	Port      uint16
	PublicKey rsa.PublicKey
}

func NewIDManagerData(IP net.IP, Port uint16, PublicKey rsa.PublicKey) IDManagerData {
	return IDManagerData{
		IP:        IP,
		Port:      Port,
		PublicKey: PublicKey,
	}
}
