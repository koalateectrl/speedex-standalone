package tatonnement

import (
	"math"

	"github.com/sandymule/speedex-standalone/orderbook"
)

type ControlParams struct {
	MTaxRate, MSmoothMult                uint8
	MMaxRounds                           uint32
	MMStepUp, MMStepDown, MStepSizeRadix uint8
	MStepRadix                           float64
	KStartingStepSize                    float64
	KPriceMin                            float64
	KPriceMax                            float64
	MRoundNumber                         uint32
}

func (cp *ControlParams) IncrementRound() {
	cp.MRoundNumber++
}

func (cp *ControlParams) Done() bool {
	return cp.MRoundNumber >= cp.MMaxRounds
}

func (cp *ControlParams) ImposePriceBounds(candidate float64) float64 {
	if candidate > cp.KPriceMax {
		return cp.KPriceMax
	}
	if candidate < cp.KPriceMin {
		return cp.KPriceMin
	}
	return candidate
}

func (cp *ControlParams) SetTrialPrice(curPrice float64, demand float64, stepSize float64) float64 {
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

	delta := product / cp.MStepRadix

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

func (cp *ControlParams) SetTrialPrices(curPrices map[orderbook.Asset]float64, demands orderbook.SupplyDemand, stepSize float64) map[orderbook.Asset]float64 {
	// set prices for all assets
	pricesOut := make(map[orderbook.Asset]float64)
	for asset, curPrice := range curPrices {
		pricesOut[asset] = cp.SetTrialPrice(curPrice, demands.GetDelta(asset), stepSize)
	}

	return pricesOut
}
