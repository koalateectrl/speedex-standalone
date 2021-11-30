package tatonnement

import (
	"fmt"
	"math"

	"github.com/sandymule/speedex-standalone/pkg/assets"
	"github.com/sandymule/speedex-standalone/pkg/orderbook"
)

type TatonnementOracle struct {
	MOrderbookManager orderbook.OrderbookManager
}

func (to *TatonnementOracle) ComputePrices(params TatonnementControlParams, prices map[assets.Asset]float64, printFrequency uint32) {
	controlParams := NewTatonnementControlParamsWrapper(&params)
	fmt.Println(to)
	baselineDemand := to.MOrderbookManager.DemandQuery(prices, controlParams.GetSmoothMult())
	baselineObjective := baselineDemand.GetObjective()
	fmt.Println(baselineDemand)
	fmt.Println(baselineObjective)

	stepSize := controlParams.KStartingStepSize

	for !controlParams.Done() {
		controlParams.IncrementRound()
		fmt.Println(prices)
		trialPrices := controlParams.SetTrialPrices(prices, *baselineDemand, stepSize)
		trialDemand := to.MOrderbookManager.DemandQuery(trialPrices, controlParams.GetSmoothMult())
		trialObjective := trialDemand.GetObjective()

		if trialObjective.Value <= baselineObjective.Value || stepSize < controlParams.KMinStepSize {
			prices = trialPrices
			baselineDemand = trialDemand
			baselineObjective = trialObjective
			stepSize = controlParams.StepUp(math.Max(stepSize, controlParams.KMinStepSize))
		} else {
			stepSize = controlParams.StepDown(stepSize)
		}

		if printFrequency > 0 && controlParams.GetRoundNumber()%printFrequency == 0 {
			fmt.Printf("TATONNEMENT STEP: step size: %f round number: %d\n", stepSize, controlParams.MRoundNumber)
			for asset, price := range prices {
				demand := baselineDemand.GetDelta(asset)
				fmt.Printf("TATONNEMENT: %s, Price: %f, Demand: %f\n", string(asset), price, demand)
			}
		}
	}
}
