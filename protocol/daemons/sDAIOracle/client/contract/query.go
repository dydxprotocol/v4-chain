package store

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
)

func QueryDaiConversionRate(client *ethclient.Client) (string, string, error) {
	// Create an instance of the contract
	instance, err := NewStore(types.MakerContractAddress, client)
	if err != nil {
		return "", "", err
	}

	// Query the chi variable
	sDAIExchangeRate, err := instance.Chi(nil)
	if err != nil {
		return "", "", err
	}

	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return "", "", err
	}

	return sDAIExchangeRate.String(), header.Number.String(), nil
}
