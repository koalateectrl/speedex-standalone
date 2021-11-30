package orderbook

import (
	"math"
)

type PriceCompStats struct {
	SellPrice                          float64
	CumulativeOfferedForSale           float64
	CumulativeOfferedForSaleTimesPrice float64
	Txid                               uint64
}

type Orderbook struct {
	MPrecomputedTatonnementData []PriceCompStats
}

func (ob *Orderbook) GetPriceCompStats(sellPrice float64, buyPrice float64) PriceCompStats {
	// TODO insert condition when vector is empty to return zero stats
	var start uint8 = 0
	var end uint8 = uint8(len(ob.MPrecomputedTatonnementData) - 1)

	if ob.MPrecomputedTatonnementData[end].SellPrice <= sellPrice/buyPrice {
		return ob.MPrecomputedTatonnementData[end]
	}

	if sellPrice/buyPrice <= ob.MPrecomputedTatonnementData[start].SellPrice {
		return PriceCompStats{SellPrice: 0, CumulativeOfferedForSale: 0, CumulativeOfferedForSaleTimesPrice: 0}
	}

	// binary search
	for {
		var mid uint8 = (start + end) / 2
		if start == end {
			return ob.MPrecomputedTatonnementData[start-1]
		}
		if ob.MPrecomputedTatonnementData[mid].SellPrice <= sellPrice/buyPrice {
			start = mid + 1
		} else {
			end = mid
		}
	}
}

func (ob *Orderbook) ApplySmoothMult(sellPrice float64, smoothMult uint8) float64 {
	if smoothMult == 0 {
		return sellPrice
	}

	return sellPrice - sellPrice*math.Pow(2, -float64(smoothMult))
}

func (ob *Orderbook) CumulativeOfferedForSaleTimesPrice(sellPrice float64, buyPrice float64, smoothMult uint8) float64 {
	// TODO double check math for the partial Amounts (test with transaction set with small sellprice differences)
	fullExecSellPrice := ob.ApplySmoothMult(sellPrice, smoothMult)
	partialExecSellPrice := sellPrice

	fullExecStats := ob.GetPriceCompStats(fullExecSellPrice, buyPrice)
	partialExecStats := ob.GetPriceCompStats(partialExecSellPrice, buyPrice)

	fullExecEndow := fullExecStats.CumulativeOfferedForSale
	partialExecEndow := partialExecStats.CumulativeOfferedForSale - fullExecEndow

	fullExecEndowTimesPrice := fullExecStats.CumulativeOfferedForSaleTimesPrice
	partialExecEndowTimesPrice := partialExecStats.CumulativeOfferedForSaleTimesPrice - fullExecEndowTimesPrice

	partialAmountTimesSellPrice := partialExecEndow * sellPrice

	partialAmountTimesMinPriceTimesBuyPrice := partialExecEndowTimesPrice * buyPrice

	partialAmountSum := partialAmountTimesSellPrice - partialAmountTimesMinPriceTimesBuyPrice

	valueSoldPartialExec := partialAmountSum * math.Pow(2, float64(smoothMult))

	valueSoldFullExec := fullExecEndow * sellPrice

	return valueSoldFullExec + valueSoldPartialExec
}
