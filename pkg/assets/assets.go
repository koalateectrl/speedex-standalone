package assets

type Asset string

type AssetPair struct {
	buying  Asset
	selling Asset
}

/*
type Price struct {
	N float64 // numerator
	D float64 // denominator
} */

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
