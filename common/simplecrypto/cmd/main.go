package main

import (
	"eaglechat/common/simplecrypto/rsa"
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
	os.WriteFile("public_key.pem", pkBytes, 0644)

	skBytes := sk.ToBytes()
	os.WriteFile("private_key.pem", skBytes, 0600)
}
