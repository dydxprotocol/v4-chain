package constants

import (
	clobtypes "github.com/dydxprotocol/v4-chain/protocol/x/clob/types"
)

var (
	// Liquidation Orders.
	LiquidationOrder_Carl_Num0_Clob0_Buy3_Price50_BTC = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		3,
		50,
	)
	LiquidationOrder_Bob_Num0_Clob0_Buy100_Price20_BTC = *clobtypes.NewLiquidationOrder(
		Bob_Num0,
		ClobPair_Btc,
		true,
		100,
		20,
	)
	LiquidationOrder_Bob_Num0_Clob0_Buy25_Price30_BTC = *clobtypes.NewLiquidationOrder(
		Bob_Num0,
		ClobPair_Btc,
		true,
		25,
		30,
	)
	LiquidationOrder_Alice_Num0_Clob0_Sell20_Price25_BTC = *clobtypes.NewLiquidationOrder(
		Alice_Num0,
		ClobPair_Btc,
		false,
		20,
		25,
	)
	LiquidationOrder_Bob_Num0_Clob1_Sell70_Price10_ETH = *clobtypes.NewLiquidationOrder(
		Bob_Num0,
		ClobPair_Eth,
		false,
		70,
		10,
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy100BTC_Price101 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		10_000_000_000, // 100 BTC
		101_000_000,    // $101
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy01BTC_Price50000 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		10_000_000,
		50_000_000_000,
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50000 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		100_000_000,
		50_000_000_000,
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50500 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		100_000_000,
		50_500_000_000,
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price50501_01 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		100_000_000,
		50_501_010_000,
	)
	LiquidationOrder_Carl_Num0_Clob0_Buy1BTC_Price60000 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Btc,
		true,
		100_000_000,
		60_000_000_000,
	)
	LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3000 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Eth,
		true,
		1_000_000_000,
		3_000_000_000,
	)
	LiquidationOrder_Carl_Num0_Clob1_Buy1ETH_Price3030 = *clobtypes.NewLiquidationOrder(
		Carl_Num0,
		ClobPair_Eth,
		true,
		1_000_000_000,
		3_030_000_000,
	)
	LiquidationOrder_Dave_Num0_Clob0_Buy100BTC_Price102 = *clobtypes.NewLiquidationOrder(
		Dave_Num0,
		ClobPair_Btc,
		true,
		10_000_000_000, // 100 BTC
		102_000_000,    // $102
	)
	LiquidationOrder_Dave_Num0_Clob0_Sell100BTC_Price98 = *clobtypes.NewLiquidationOrder(
		Dave_Num0,
		ClobPair_Btc,
		false,
		10_000_000_000, // 100 BTC
		98_000_000,     // $98
	)
	LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price49500 = *clobtypes.NewLiquidationOrder(
		Dave_Num0,
		ClobPair_Btc,
		false,
		100_000_000,
		49_500_000_000,
	)
	LiquidationOrder_Dave_Num0_Clob0_Sell1BTC_Price50000 = *clobtypes.NewLiquidationOrder(
		Dave_Num0,
		ClobPair_Btc,
		false,
		100_000_000,
		50_000_000_000,
	)
	LiquidationOrder_Dave_Num1_Clob0_Buy100BTC_Price101 = *clobtypes.NewLiquidationOrder(
		Dave_Num1,
		ClobPair_Btc,
		true,
		10_000_000_000, // 100 BTC
		101_000_000,    // $101
	)
	LiquidationOrder_Dave_Num1_Clob0_Buy100BTC_Price102 = *clobtypes.NewLiquidationOrder(
		Dave_Num1,
		ClobPair_Btc,
		true,
		10_000_000_000, // 100 BTC
		102_000_000,    // $102
	)
	LiquidationOrder_Dave_Num1_Clob0_Sell01BTC_Price50000 = *clobtypes.NewLiquidationOrder(
		Dave_Num1,
		ClobPair_Btc,
		false,
		10_000_000,
		50_000_000_000,
	)
)
