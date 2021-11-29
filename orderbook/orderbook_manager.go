package orderbook

type OrderbookManager struct {
	MOrderbooks map[AssetPair]Orderbook
}

func (obm *OrderbookManager) DemandQuery(prices map[Asset]float64, smoothMult uint8) *SupplyDemand {
	supplyDemand := &SupplyDemand{MSupplyDemand: make(map[Asset]*SupplyDemandPair)}

	for assetPair, ob := range obm.MOrderbooks {
		sellPrice := prices[assetPair.selling]
		buyPrice := prices[assetPair.buying]

		tradeAmount := ob.CumulativeOfferedForSaleTimesPrice(sellPrice, buyPrice, smoothMult)
		supplyDemand.AddSupplyDemandPair(assetPair, tradeAmount)
	}

	return supplyDemand
}
