package entities

import (
	"net"

	"eaglechat/common/simplecrypto/rsa"
)

const DefaultIDManagerPort uint16 = 8080

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
