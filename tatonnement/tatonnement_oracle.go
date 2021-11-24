package tatonnement

import (
	"fmt"

	"github.com/sandymule/speedex-standalone/orderbook"
)

type TatonnementOracle struct {
	Params            ControlParams
	MOrderbookManager orderbook.OrderbookManager
}

func (to *TatonnementOracle) ComputePrices(params ControlParams, prices map[orderbook.Asset]float64) {
	to.Params = params
	fmt.Println(to)
	baselineDemand := to.MOrderbookManager.DemandQuery(prices)
	fmt.Println(baselineDemand)
}
