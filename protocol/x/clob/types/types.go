package types

import "math/big"

type FillForProcess interface {
	TakerAddr() string
	TakerFeeQuoteQuantums() *big.Int
	MakerAddr() string
	MakerFeeQuoteQuantums() *big.Int
	FillQuoteQuantums() *big.Int
	FillSourceId() uint32
	MonthlyRollingTakerVolumeQuantums() uint64
}

type PerpetualFillForProcess struct {
	takerAddr                         string
	takerFeeQuoteQuantums             *big.Int
	makerAddr                         string
	makerFeeQuoteQuantums             *big.Int
	fillQuoteQuantums                 *big.Int
	perpetualId                       uint32
	monthlyRollingTakerVolumeQuantums uint64
}

func (perpetualFillForProcess PerpetualFillForProcess) TakerAddr() string {
	return perpetualFillForProcess.takerAddr
}

func (perpetualFillForProcess PerpetualFillForProcess) TakerFeeQuoteQuantums() *big.Int {
	return perpetualFillForProcess.takerFeeQuoteQuantums
}

func (perpetualFillForProcess PerpetualFillForProcess) MakerAddr() string {
	return perpetualFillForProcess.makerAddr
}

func (perpetualFillForProcess PerpetualFillForProcess) MakerFeeQuoteQuantums() *big.Int {
	return perpetualFillForProcess.makerFeeQuoteQuantums
}

func (perpetualFillForProcess PerpetualFillForProcess) FillQuoteQuantums() *big.Int {
	return perpetualFillForProcess.fillQuoteQuantums
}

func (perpetualFillForProcess PerpetualFillForProcess) FillSourceId() uint32 {
	return perpetualFillForProcess.perpetualId
}

func (perpetualFillForProcess PerpetualFillForProcess) MonthlyRollingTakerVolumeQuantums() uint64 {
	return perpetualFillForProcess.monthlyRollingTakerVolumeQuantums
}

func CreatePerpetualFillForProcess(
	takerAddr string,
	takerFeeQuoteQuantums *big.Int,
	makerAddr string,
	makerFeeQuoteQuantums *big.Int,
	fillQuoteQuantums *big.Int,
	perpetualId uint32,
	monthlyRollingTakerVolumeQuantums uint64,
) PerpetualFillForProcess {
	return PerpetualFillForProcess{
		takerAddr:                         takerAddr,
		takerFeeQuoteQuantums:             takerFeeQuoteQuantums,
		makerAddr:                         makerAddr,
		makerFeeQuoteQuantums:             makerFeeQuoteQuantums,
		fillQuoteQuantums:                 fillQuoteQuantums,
		perpetualId:                       perpetualId,
		monthlyRollingTakerVolumeQuantums: monthlyRollingTakerVolumeQuantums,
	}
}
