package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Pack store(42)
//     ↓
// Create Ethereum transaction
//     ↓
// Sign transaction
//     ↓
// Send transaction to Sepolia
//     ↓
// Wait for confirmation
//     ↓
// Read retrieve()
//     ↓
// 42

const contractABI = `[
	{
		"inputs": [
			{
				"internalType": "uint256",
				"name": "_favoriteNumber",
				"type": "uint256"
			}
		],
		"name": "store",
		"outputs": [],
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

func main() {

	ctx := context.Background()

	// Connect to Sepolia
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/acUFVE2BQoKhWJQR-d5Fx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Sepolia")

	// Contract address
	contractAddress := common.HexToAddress(
		"0x41b727f6e4e689589F3549f8FcB0FC2c4299af41",
	)

	// Load private key... You shouldn't display your private keys this way.
	privateKeyHex := "7e9458d590a5bc4381576c326380e6f66f7983523817dbeea264946e45391a6c"

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Get public key
	publicKey := privateKey.Public()

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Unable to get public key")
	}

	// Get wallet address
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Println("From:", fromAddress.Hex())

	// Parse ABI
	parsedABI, err := abi.JSON(
		strings.NewReader(contractABI),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Value we want to store
	favoriteNumber := big.NewInt(42)

	// Encode store(42)
	data, err := parsedABI.Pack(
		"store",
		favoriteNumber,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get account nonce
	nonce, err := client.PendingNonceAt(
		ctx,
		fromAddress,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get chain ID
	chainID, err := client.ChainID(
		ctx,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Estimate gas
	gasLimit, err := client.EstimateGas(
		ctx,
		ethereum.CallMsg{
			From: fromAddress,
			To:   &contractAddress,
			Data: data,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get suggested gas price
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Create transaction
	tx := types.NewTransaction(
		nonce,
		contractAddress,
		big.NewInt(0),
		gasLimit,
		gasPrice,
		data,
	)

	// Sign transaction
	signedTx, err := types.SignTx(
		tx,
		types.NewEIP155Signer(chainID),
		privateKey,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Send transaction
	err = client.SendTransaction(
		ctx,
		signedTx,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transaction sent!")
	fmt.Println("Transaction hash:", signedTx.Hash().Hex())
}

// 0x97697a665b9315D125a5C389B28aBf99209713C0
