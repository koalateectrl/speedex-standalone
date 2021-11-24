package orderbook

import "fmt"

type OrderbookManager struct {
	MOrderbooks map[AssetPair]Orderbook
}

func (obm *OrderbookManager) DemandQuery(prices map[Asset]float64) *SupplyDemand {
	supplyDemand := &SupplyDemand{MSupplyDemand: make(map[Asset]*SupplyDemandPair)}

	ethSupplyDemandPair := SupplyDemandPair{1, 1}
	supplyDemand.MSupplyDemand["ETH"] = &ethSupplyDemandPair
	usdtSupplyDemandPair := SupplyDemandPair{2, 2}
	supplyDemand.MSupplyDemand["USDT"] = &usdtSupplyDemandPair

	for assetPair, ob := range obm.MOrderbooks {
		fmt.Println(assetPair)
		sellPrice := prices[assetPair.selling]
		buyPrice := prices[assetPair.buying]

		tradeAmount := ob.CumulativeOfferedForSaleTimesPrice(sellPrice, buyPrice)
		fmt.Println(tradeAmount)
		supplyDemand.AddSupplyDemandPair(assetPair, tradeAmount)
	}

	return supplyDemand
}
