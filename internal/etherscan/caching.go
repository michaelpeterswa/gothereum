package etherscan

import (
	"encoding/json"
	"os"
	"time"
)

const tempFile = "/tmp/gothereum.json"

type GothereumCache struct {
	Timestamp                    time.Time `json:"timestamp"`
	EthereumToUnitedStatesDollar float64   `json:"ethusd"`
}

func removeLastPrice() error {
	return os.Remove(tempFile)
}

func saveLastPrice(price float64) error {
	removeLastPrice()

	file, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(GothereumCache{Timestamp: time.Now(), EthereumToUnitedStatesDollar: price})
}

func loadLastPrice() (*GothereumCache, error) {
	var gc GothereumCache

	file, err := os.Open(tempFile)
	// safely skip no file
	if err != nil {
		return &GothereumCache{Timestamp: time.Now(), EthereumToUnitedStatesDollar: 0}, nil
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&gc)
	if err != nil {
		return nil, err
	}

	return &gc, nil
}
