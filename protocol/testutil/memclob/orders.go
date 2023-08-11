package memclob

// OrderSideHumanReadable is a testing observability function that translates the boolean side of an order to a human
// readable string. For example, passing in `true` represents a buy order, so this function will return `"BUY"`,
// whereas passing in `false` will return `"SELL"`.
func OrderSideHumanReadable(isBuy bool) string {
	if isBuy {
		return "BUY"
	} else {
		return "SELL"
	}
}
