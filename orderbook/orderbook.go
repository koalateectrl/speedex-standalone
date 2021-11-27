package orderbook

import (
	"sort"
)

type PriceCompStats struct {
	SellPrice                float64
	CumulativeOfferedForSale float64
	Txid                     uint64
}

type Orderbook struct {
	MPrecomputedTatonnementData []PriceCompStats
}

func (ob *Orderbook) CumulativeOfferedForSaleTimesPrice(sellPrice float64, buyPrice float64) float64 {
	//TODO code for partial/full sells here
	//TODO add code to price weight the offers
	p := sellPrice / buyPrice
	pos := sort.Search(len(ob.MPrecomputedTatonnementData), func(i int) bool { return ob.MPrecomputedTatonnementData[i].SellPrice > p })
	if pos == 0 { // no limit price satisfies current price
		return 0
	}
	return ob.MPrecomputedTatonnementData[pos-1].CumulativeOfferedForSale * sellPrice
}
