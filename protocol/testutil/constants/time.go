package constants

import (
	"time"

	"github.com/dydxprotocol/v4/daemons/pricefeed"
)

var (
	// Time
	TimeZero            = time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	TimeT               = time.Unix(1650000000, 0) // 2022-04-14 22:20:00 -0700 PDT
	TimeTPlus1          = TimeT.Add(time.Duration(1))
	TimeTMinusThreshold = TimeT.Add(-pricefeed.MaxPriceAge).Add(-time.Duration(1))
	TimeTPlusThreshold  = TimeT.Add(pricefeed.MaxPriceAge).Add(time.Duration(1))
	Time_21st_Feb_2021  = time.Date(2021, time.Month(2), 21, 0, 0, 0, 0, time.UTC)
)
