package orderbook

import "strconv"

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
	if _, found := sd.MSupplyDemand[sell]; found {
		sd.MSupplyDemand[sell].Supply += amount
	} else {
		sd.MSupplyDemand[sell] = &SupplyDemandPair{amount, 0}
	}

	if _, found := sd.MSupplyDemand[buy]; found {
		sd.MSupplyDemand[buy].Demand += amount
	} else {
		sd.MSupplyDemand[buy] = &SupplyDemandPair{0, amount}
	}
}

func (sd *SupplyDemand) String() string {
	retStr := "{"
	for key, val := range sd.MSupplyDemand {
		retStr += string(key) + ": ("
		retStr += strconv.FormatFloat(val.Supply, 'f', -1, 64) + "," + strconv.FormatFloat(val.Demand, 'f', -1, 64) + "); "
	}
	retStr += "}"
	return retStr
}
