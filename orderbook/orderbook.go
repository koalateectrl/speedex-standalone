package orderbook

type Asset string

type AssetPair struct {
	buying  Asset
	selling Asset
}

type AssetTx struct {
	Type       string
	Amount     uint64
	LimitPrice uint64
}

type Transactions struct {
	Transactions []Transaction `json:"txs"`
}

type Transaction struct {
	Txid      uint64  `json:"txid"`
	BuyAsset  AssetTx `json:"buyasset"`
	SellAsset AssetTx `json:"sellasset"`
}

type IOCOrderbook struct {
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
