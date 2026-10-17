package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const retrieveABI = `[
		{
			"inputs": [],
			"name": "retrieve",
			"outputs": [
				{
					"internalType": "uint256",
					"name": "",
					"type": "uint256"
				}
			],
			"stateMutability": "view",
			"type": "function"
		}
	]`

func main() {

	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/acUFVE2BQoKhWJQR-d5Fx")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to Ethereum")

	// Contract address
	contractAddress := common.HexToAddress(
		"0x41b727f6e4e689589f3549f8fcb0fc2c4299af41",
	)

	// Parse ABI
	parsedABI, err := abi.JSON(
		strings.NewReader(retrieveABI),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Encode retrieve()
	data, err := parsedABI.Pack("retrieve")
	if err != nil {
		log.Fatal(err)
	}

	// Call contract
	result, err := client.CallContract(
		context.Background(),
		ethereum.CallMsg{
			To:   &contractAddress,
			Data: data,
		},
		nil,
	)

	if err != nil {
		log.Fatal(err)
	}

	// Decode result
	output, err := parsedABI.Unpack("retrieve", result)
	if err != nil {
		log.Fatal(err)
	}

	favoriteNumber := output[0].(*big.Int)

	fmt.Println("Favorite number:", favoriteNumber)
}
