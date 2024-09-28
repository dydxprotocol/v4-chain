package store

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sdaioracle/client/types"
)

func QueryDaiConversionRate(client *ethclient.Client) (string, error) {
	// Create an instance of the contract
	instance, err := NewStore(types.MakerContractAddress, client)
	if err != nil {
		return "", err
	}

	// Query the chi variable
	sDAIExchangeRate, err := instance.Chi(nil)
	if err != nil {
		return "", err
	}

	return sDAIExchangeRate.String(), nil
}

func QueryDaiConversionRateWithRetries(client *ethclient.Client, maxRetries int) (string, error) {

	for i := 0; i < maxRetries; i++ {
		rate, err := QueryDaiConversionRate(client)
		if err == nil {
			return rate, nil
		}
		time.Sleep(time.Second)
	}
	return "", errors.New("failed to query DAI conversion rate")
}
