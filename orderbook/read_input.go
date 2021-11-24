package orderbook

type Transactions struct {
	Transactions []Transaction `json:"txs"`
}

type Transaction struct {
	Txid      uint64  `json:"txid"`
	BuyAsset  AssetTx `json:"buyasset"`
	SellAsset AssetTx `json:"sellasset"`
}
