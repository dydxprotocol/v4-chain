package store

import (
	"context"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

func QueryDaiConversionRateForPastBlocks(client *ethclient.Client, blocks int64, maxRetries int) ([]string, []string, error) {
	var rates []string
	var blockNumbers []string

	// Get the latest block number
	latestHeader, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	latestBlockNumber := latestHeader.Number.Int64()

	// Create an instance of the contract
	instance, err := NewStore(types.MakerContractAddress, client)
	if err != nil {
		return nil, nil, err
	}

	for i := int64(0); i < blocks; i++ {
		blockNumber := latestBlockNumber - i
		var sDAIExchangeRate *big.Int

		for retry := 0; retry < maxRetries; retry++ {
			// Query the chi variable for the specific block
			sDAIExchangeRate, err = instance.Chi(&bind.CallOpts{
				BlockNumber: big.NewInt(blockNumber),
			})
			if err == nil {
				break
			}
			if retry == maxRetries-1 || !strings.Contains(err.Error(), "capacity") {
				return nil, nil, err
			}

			time.Sleep(time.Second * 1)
		}

		rates = append(rates, sDAIExchangeRate.String())
		blockNumbers = append(blockNumbers, big.NewInt(blockNumber).String())
	}

	return rates, blockNumbers, nil
}
