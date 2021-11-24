package orderbook

type OrderbookManager struct {
	MOrderbooks map[AssetPair]Orderbook
}
