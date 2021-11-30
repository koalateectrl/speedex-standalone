package tatonnement

import (
	"math"

	"github.com/sandymule/speedex-standalone/pkg/assets"
	"github.com/sandymule/speedex-standalone/pkg/demandutils"
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

func (tcpw *TatonnementControlParamsWrapper) GetRoundNumber() uint32 {
	return tcpw.MRoundNumber
}

func (tcpw *TatonnementControlParamsWrapper) GetSmoothMult() uint8 {
	return tcpw.MParams.MSmoothMult
}

func (tcpw *TatonnementControlParamsWrapper) IncrementRound() {
	tcpw.MRoundNumber++
}

func (tcpw *TatonnementControlParamsWrapper) Done() bool {
	return tcpw.MRoundNumber >= tcpw.MParams.MMaxRounds
}

func (tcpw *TatonnementControlParamsWrapper) ImposePriceBounds(candidate float64) float64 {
	if candidate > tcpw.KPriceMax {
		return tcpw.KPriceMax
	}
	if candidate < tcpw.KPriceMin {
		return tcpw.KPriceMin
	}
	return candidate
}

func (tcpw *TatonnementControlParamsWrapper) StepUp(step float64) float64 {
	out := step * float64(tcpw.MParams.MStepUp) * math.Pow(2, -float64(tcpw.MParams.MStepSizeRadix))
	return out
}

func (tcpw *TatonnementControlParamsWrapper) StepDown(step float64) float64 {
	out := step * float64(tcpw.MParams.MStepDown) * math.Pow(2, -float64(tcpw.MParams.MStepSizeRadix))
	return out
}

func (tcpw *TatonnementControlParamsWrapper) SetTrialPrice(curPrice float64, demand float64, stepSize float64) float64 {
	// set price for one asset
	stepTimesOldPrice := curPrice * stepSize

	var sign int8 = 1
	var unsignedDemand float64 = demand
	if demand <= 0 {
		sign = -1
		unsignedDemand = -demand
	}

	product := stepTimesOldPrice * unsignedDemand

	delta := product * math.Pow(2, -float64(tcpw.MParams.MStepRadix))

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
	return tcpw.ImposePriceBounds(candidateOut)
}

func (tcpw *TatonnementControlParamsWrapper) SetTrialPrices(curPrices map[assets.Asset]float64, demands demandutils.SupplyDemand, stepSize float64) map[assets.Asset]float64 {
	// set prices for all assets
	pricesOut := make(map[assets.Asset]float64)
	for asset, curPrice := range curPrices {
		pricesOut[asset] = tcpw.SetTrialPrice(curPrice, demands.GetDelta(asset), stepSize)
	}

	return pricesOut
}
