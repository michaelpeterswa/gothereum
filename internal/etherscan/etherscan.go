package etherscan

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
)

const accountRequestURL = "https://api.etherscan.io/api?module=account&action=balance&address=%s&tag=latest&apikey=%s"
const ethereumPriceURL = "https://api.etherscan.io/api?module=stats&action=ethprice&apikey=%s"

type EtherscanClient struct {
	client *http.Client
	apiKey string
}

type AccountBalanceResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

type EthereumPriceResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Result  EthereumPriceResult `json:"result"`
}

type EthereumPriceResult struct {
	EthereumToBitcoin                     string `json:"ethbtc"`
	EthereumToBitcoinTimestamp            string `json:"ethbtc_timestamp"`
	EthereumToUnitedStatesDollar          string `json:"ethusd"`
	EthereumToUnitedStatesDollarTimestamp string `json:"ethusd_timestamp"`
}

func NewEtherscanClient(apiKey string) *EtherscanClient {
	return &EtherscanClient{
		client: &http.Client{Timeout: 10 * time.Second},
		apiKey: apiKey,
	}
}

func (e *EtherscanClient) getAccountBalance(accountAddress string) (float64, error) {
	var abr AccountBalanceResponse

	req, err := http.NewRequest("GET", fmt.Sprintf(accountRequestURL, accountAddress, e.apiKey), nil)
	if err != nil {
		return 0, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&abr)
	if err != nil {
		return 0, err
	}

	wei, err := strconv.ParseInt(abr.Result, 10, 64)
	if err != nil {
		return 0, err
	}

	return weiToEther(wei), nil
}

func (e *EtherscanClient) getEthereumPrices() (float64, float64, error) {
	var epr EthereumPriceResponse

	req, err := http.NewRequest("GET", fmt.Sprintf(ethereumPriceURL, e.apiKey), nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&epr)
	if err != nil {
		return 0, 0, err
	}

	gc, err := loadLastPrice()
	if err != nil {
		return 0, 0, err
	}

	ethToUsd, err := strconv.ParseFloat(epr.Result.EthereumToUnitedStatesDollar, 64)
	if err != nil {
		return 0, 0, err
	}

	err = saveLastPrice(ethToUsd)
	if err != nil {
		return 0, 0, err
	}

	return ethToUsd, gc.EthereumToUnitedStatesDollar, nil
}

func weiToEther(wei int64) float64 {
	return float64(wei) / math.Pow10(18)
}

func (e *EtherscanClient) GetAccountBalanceUSD(accountAddress string) (float64, string, error) {
	bal, err := e.getAccountBalance(accountAddress)
	if err != nil {
		return 0, "", err
	}

	price, historicPrice, err := e.getEthereumPrices()
	if err != nil {
		return 0, "", err
	}

	if price == historicPrice {
		return bal * price, "equal", nil
	} else if price > historicPrice {
		return bal * price, "more", nil
	}
	return bal * price, "less", nil
}
