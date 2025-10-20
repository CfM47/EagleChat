package main

import (
	"eaglechat/common/ezcrypto/rsa"
	"os"
)

func main() {
	sk, pk, err := rsa.GenerateKeyPair()
	if err != nil {
		panic(err)
	}

	pkBytes, err := pk.ToBytes()
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("public_key.pem", pkBytes, 0644)
	if err != nil {
		panic(err)
	}

	skBytes := sk.ToBytes()
	err = os.WriteFile("private_key.pem", skBytes, 0600)
	if err != nil {
		panic(err)
	}
}
