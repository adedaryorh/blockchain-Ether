package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// Connect to Ethereum
	//https://ethereum-sepolia-rpc.publicnode.com"
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/acUFVE2BQoKhWJQR-d5Fx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Sepolia")

	// Load Private Key
	privateKeyHex := "7e9458d590a5bc4381576c326380e6f66f7983523817dbeea264946e45391a6c"

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Derive Sender Address
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot cast public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("Sender:", fromAddress.Hex())

	// Recipient Address
	toAddress := common.HexToAddress("0x3BD82d57f63ef1CDB5D177A6B149312B4d80D390")

	// Next Nonce
	nonce, err := client.PendingNonceAt(
		context.Background(),
		fromAddress,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Amount (0.01 ETH)
	value := big.NewInt(10000000000000000)

	// Gas
	gasLimit := uint64(21000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("suggested gas: ", gasPrice)

	// Create Transaction
	tx := types.NewTransaction(
		nonce,
		toAddress,
		value,
		gasLimit,
		gasPrice,
		nil,
	)

	// Chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("sepolia chainID: ", chainID)

	// Sign Transaction
	signedTx, err := types.SignTx(
		tx,
		types.LatestSignerForChainID(chainID),
		privateKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Broadcast Transaction
	err = client.SendTransaction(
		context.Background(),
		signedTx,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction sent successfully!")
	fmt.Println("Transaction Hash:", signedTx.Hash().Hex())
}
