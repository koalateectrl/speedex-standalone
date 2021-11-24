package orderbook

import "fmt"

type Asset string

type AssetPair struct {
	buying  Asset
	selling Asset
}

type SupplyDemandPair struct {
	Supply float64
	Demand float64
}

type SupplyDemand struct {
	MSupplyDemand map[Asset]*SupplyDemandPair
}

func (ap *AssetPair) Buying() Asset {
	return ap.buying
}

func (ap *AssetPair) Selling() Asset {
	return ap.selling
}

func (ap *AssetPair) SetBuying(newBuy Asset) {
	ap.buying = newBuy
}

func (ap *AssetPair) SetSelling(newSell Asset) {
	ap.selling = newSell
}

func (ap *AssetPair) String() string {
	return "(" + string(ap.buying) + " / " + string(ap.selling) + ")"
}

func (sd *SupplyDemand) AddSupplyDemandPair(tradingPair AssetPair, amount float64) {
	sd.AddSupplyDemand(tradingPair.selling, tradingPair.buying, amount)
}

func (sd *SupplyDemand) AddSupplyDemand(sell Asset, buy Asset, amount float64) {
	fmt.Println(sd.MSupplyDemand[sell])
	sd.MSupplyDemand[sell].Supply += amount
	sd.MSupplyDemand[buy].Demand += amount
}
