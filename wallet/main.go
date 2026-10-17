package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	privateKeyHex := "7e9458d590a5bc4381576c326380e6f66f7983523817dbeea264946e45391a6c"

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("Wallet Address:", address.Hex())
	log.Println("metamask: ", "0x97697a665b9315D125a5C389B28aBf99209713C0")
}
