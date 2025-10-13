package entities

import (
	"net"
)

const DefaultIDManagerPort uint16 = 8080

type IDManagerData struct {
	IP   net.IP
	Port uint16

	// TODO: Backlog
	//
	// PublicKey rsa.PublicKey
}

// TODO: Backlog
//
// func NewIDManagerData(IP net.IP, Port uint16, PublicKey rsa.PublicKey) IDManagerData {

func NewIDManagerData(IP net.IP, Port uint16) IDManagerData {
	return IDManagerData{
		IP:   IP,
		Port: Port,
		// PublicKey: PublicKey,
	}
}
