package main

import (
	"errors"
	"fmt"
	"net"
	"os"

	"github.com/michaelpeterswa/gothereum/internal/etherscan"
	"github.com/michaelpeterswa/gothereum/internal/textutils"
)

func main() {
	etherscanKey := os.Getenv("ETHERSCAN_API_KEY")
	if etherscanKey == "" {
		fmt.Println("$ETHERSCAN_API_KEY not found...")
		os.Exit(1)
	}

	ethereumAddress := os.Getenv("ETHEREUM_ADDRESS")
	if ethereumAddress == "" {
		fmt.Println("$ETHEREUM_ADDRESS not found...")
		os.Exit(1)
	}

	etherscanClient := etherscan.NewEtherscanClient(etherscanKey)

	balance, status, err := etherscanClient.GetAccountBalanceUSD(ethereumAddress)

	var dnsError *net.DNSError
	if errors.As(err, &dnsError) {
		os.Exit(0)
	}

	fmt.Println(textutils.InsertSymbolAndColor(fmt.Sprintf("%.2f USD", balance), status))

}
