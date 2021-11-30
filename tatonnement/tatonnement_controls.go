package tatonnement

import (
	"math"

	"github.com/sandymule/speedex-standalone/assets"
	"github.com/sandymule/speedex-standalone/demandutils"
)

type TatonnementControlParams struct {
	MSmoothMult                        uint8
	MMaxRounds                         uint32
	MStepUp, MStepDown, MStepSizeRadix uint8
	MStepRadix                         float64
}

type TatonnementControlParamsWrapper struct {
	MParams           TatonnementControlParams
	MRoundNumber      uint32
	KPriceMin         float64
	KPriceMax         float64
	KMinStepSize      float64
	KStartingStepSize float64
}

func NewTatonnementControlParamsWrapper(tcp *TatonnementControlParams) *TatonnementControlParamsWrapper {
	tcpw := new(TatonnementControlParamsWrapper)
	tcpw.MParams = *tcp
	tcpw.KPriceMin = 1
	tcpw.KPriceMax = math.MaxInt64 * math.Pow(2, -float64(tcp.MSmoothMult+1))
	tcpw.KMinStepSize = math.Pow(2, float64(tcp.MStepSizeRadix+1))
	tcpw.KStartingStepSize = tcpw.KMinStepSize
	return tcpw
}

func (cp *TatonnementControlParamsWrapper) IncrementRound() {
	cp.MRoundNumber++
}

func (cp *TatonnementControlParamsWrapper) Done() bool {
	return cp.MRoundNumber >= cp.MParams.MMaxRounds
}

func (cp *TatonnementControlParamsWrapper) ImposePriceBounds(candidate float64) float64 {
	if candidate > cp.KPriceMax {
		return cp.KPriceMax
	}
	if candidate < cp.KPriceMin {
		return cp.KPriceMin
	}
	return candidate
}

func (cp *TatonnementControlParamsWrapper) SetTrialPrice(curPrice float64, demand float64, stepSize float64) float64 {
	// set price for one asset
	stepTimesOldPrice := curPrice * stepSize

	var sign int8
	if demand > 0 {
		sign = 1
	} else {
		sign = -1
	}

	var unsignedDemand float64
	if demand > 0 {
		unsignedDemand = demand
	} else {
		unsignedDemand = -demand
	}

	product := stepTimesOldPrice * unsignedDemand

	delta := product / cp.MParams.MStepRadix

	var candidateOut float64
	if sign > 0 {
		candidateOut = curPrice + delta
		if candidateOut < curPrice {
			candidateOut = math.MaxFloat64 // overflow
		}
	} else {
		if curPrice > delta {
			candidateOut = curPrice - delta
		} else {
			candidateOut = 0
		}
	}

	return cp.ImposePriceBounds(candidateOut)
}

func (cp *TatonnementControlParamsWrapper) SetTrialPrices(curPrices map[assets.Asset]float64, demands demandutils.SupplyDemand, stepSize float64) map[assets.Asset]float64 {
	// set prices for all assets
	pricesOut := make(map[assets.Asset]float64)
	for asset, curPrice := range curPrices {
		pricesOut[asset] = cp.SetTrialPrice(curPrice, demands.GetDelta(asset), stepSize)
	}

	return pricesOut
}

func (cp *TatonnementControlParamsWrapper) StepUp(step float64) float64 {
	out := step * float64(cp.MParams.MStepUp)
	return out
}

func (cp *TatonnementControlParamsWrapper) StepDown(step float64) float64 {
	out := step * float64(cp.MParams.MStepDown)
	return out
}
