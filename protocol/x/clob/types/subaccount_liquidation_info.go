package types

func (sli *SubaccountLiquidationInfo) HasPerpetualBeenLiquidatedForSubaccount(
	perpetualId uint32,
) bool {
	for _, liquidatedPerpetualId := range sli.PerpetualsLiquidated {
		if perpetualId == liquidatedPerpetualId {
			return true
		}
	}

	return false
}
