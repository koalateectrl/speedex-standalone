package orderbook

type PriceCompStats struct {
	SellPrice                uint64
	CumulativeOfferedForSale uint64
	BuyPrice                 uint64
	Txid                     uint64
}

type Orderbook struct {
	MPrecomputedTatonnementData []PriceCompStats
}
