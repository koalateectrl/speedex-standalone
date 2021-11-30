package tatonnement

import (
	"fmt"
	"math"

	"github.com/sandymule/speedex-standalone/assets"
	"github.com/sandymule/speedex-standalone/orderbook"
)

type TatonnementOracle struct {
	MOrderbookManager orderbook.OrderbookManager
}

func (to *TatonnementOracle) ComputePrices(params TatonnementControlParams, prices map[assets.Asset]float64, printFrequency uint32) {
	controlParams := NewTatonnementControlParamsWrapper(&params)
	fmt.Println(to)
	baselineDemand := to.MOrderbookManager.DemandQuery(prices, controlParams.MParams.MSmoothMult)
	baselineObjective := baselineDemand.GetObjective()
	fmt.Println(baselineDemand)
	fmt.Println(baselineObjective)

	stepSize := controlParams.KStartingStepSize

	for !controlParams.Done() {
		controlParams.IncrementRound()
		fmt.Println(prices)
		trialPrices := controlParams.SetTrialPrices(prices, *baselineDemand, stepSize)
		trialDemand := to.MOrderbookManager.DemandQuery(trialPrices, controlParams.MParams.MSmoothMult)
		trialObjective := trialDemand.GetObjective()

		if trialObjective.Value <= baselineObjective.Value || stepSize < controlParams.KMinStepSize {
			prices = trialPrices
			baselineDemand = trialDemand
			baselineObjective = trialObjective
			stepSize = controlParams.StepUp(math.Max(stepSize, controlParams.KMinStepSize))
		} else {
			stepSize = controlParams.StepDown(stepSize)
		}

		if printFrequency > 0 && controlParams.MRoundNumber%printFrequency == 0 {
			fmt.Printf("TATONNEMENT STEP: step size: %f round number: %d\n", stepSize, controlParams.MRoundNumber)
		}
	}
}
