package constants

import (
	"time"

	pricefeed "github.com/dydxprotocol/v4-chain/protocol/daemons/pricefeed/types"
)

var (
	// Time
	TimeZero            = time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	TimeTen             = time.Unix(10, 0)
	TimeFifteen         = time.Unix(15, 0)
	TimeTwenty          = time.Unix(20, 0)
	TimeTwentyFive      = time.Unix(25, 0)
	TimeThirty          = time.Unix(30, 0)
	TimeT               = time.Unix(1650000000, 0) // 2022-04-14 22:20:00 -0700 PDT
	TimeTMinus1         = TimeT.Add(-time.Duration(1))
	TimeTPlus1          = TimeT.Add(time.Duration(1))
	TimeTMinusThreshold = TimeT.Add(-pricefeed.MaxPriceAge).Add(-time.Duration(1))
	TimeTPlusThreshold  = TimeT.Add(pricefeed.MaxPriceAge).Add(time.Duration(1))
	Time_21st_Feb_2021  = time.Date(2021, time.Month(2), 21, 0, 0, 0, 0, time.UTC)
)
