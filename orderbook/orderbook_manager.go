package orderbook

import (
	"github.com/sandymule/speedex-standalone/assets"
	"github.com/sandymule/speedex-standalone/demandutils"
)

type OrderbookManager struct {
	MOrderbooks map[assets.AssetPair]Orderbook
}

func (obm *OrderbookManager) DemandQuery(prices map[assets.Asset]float64, smoothMult uint8) *demandutils.SupplyDemand {
	supplyDemand := &demandutils.SupplyDemand{MSupplyDemand: make(map[assets.Asset]*demandutils.SupplyDemandPair)}

	for assetPair, ob := range obm.MOrderbooks {
		sellPrice := prices[assetPair.Selling()]
		buyPrice := prices[assetPair.Buying()]
		tradeAmount := ob.CumulativeOfferedForSaleTimesPrice(sellPrice, buyPrice, smoothMult)
		supplyDemand.AddSupplyDemandPair(assetPair, tradeAmount)
	}

	return supplyDemand
}
