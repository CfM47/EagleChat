package entities

import (
	"eaglechat/common/simplecrypto/rsa"
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
