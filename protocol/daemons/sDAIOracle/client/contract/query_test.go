package store

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"

	"github.com/StreamFinance-Protocol/stream-chain/protocol/daemons/sDAIOracle/client/types"
)

func TestQueryDaiConversionRate(t *testing.T) {
	// Test with uninitialized client
	assert.Panics(t, func() {
		_, _, _ = QueryDaiConversionRate(nil)
	}, "Expected panic with uninitialized client")

	// Test with real client
	client, err := ethclient.Dial(types.ETHRPC)
	if err != nil {
		t.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	chi, blockNumber, err := QueryDaiConversionRate(client)
	assert.Nil(t, err, "Expected no error with real client")

	// Convert block number to big.Int
	blockNum := new(big.Int)
	blockNum, ok := blockNum.SetString(blockNumber, 10)
	assert.True(t, ok, "Failed to convert block number to big.Int")

	// Check block number range
	assert.True(t, blockNum.Cmp(big.NewInt(20000000)) >= 0 && blockNum.Cmp(big.NewInt(21000000)) <= 0, "Block number out of expected range")

	// Convert chi to big.Int
	sDAIExchangeRate := new(big.Int)
	sDAIExchangeRate, ok = sDAIExchangeRate.SetString(chi, 10)
	assert.True(t, ok, "Failed to convert chi to big.Int")

	// Check sDAIExchangeRate value range
	expectedsDAIExchangeRateMinValue := new(big.Int)
	expectedsDAIExchangeRateMinValue.SetString("1090000000000000000000000000", 10)
	expectedsDAIExchangeRateMaxValue := new(big.Int)
	expectedsDAIExchangeRateMaxValue.SetString("1100000000000000000000000000", 10)
	assert.True(t, sDAIExchangeRate.Cmp(expectedsDAIExchangeRateMinValue) >= 0 && sDAIExchangeRate.Cmp(expectedsDAIExchangeRateMaxValue) <= 0, "sDAIExchangeRate value out of expected range")

	// uncomment this block of code to log the results
	// log.Printf("Block Number: %s", blockNumber)
	// log.Printf("Chi Value: %s", chi)
	// panic("stop")
}
