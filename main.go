package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/sandymule/speedex-standalone/orderbook"
	"github.com/sandymule/speedex-standalone/tatonnement"
)

func createOrderbookManager(txs *orderbook.Transactions) *orderbook.OrderbookManager {
	obm := &orderbook.OrderbookManager{MOrderbooks: make(map[orderbook.AssetPair]orderbook.Orderbook)}
	for i := 0; i < len(txs.Transactions); i++ {
		var ap orderbook.AssetPair
		ap.SetBuying(orderbook.Asset(txs.Transactions[i].BuyType))
		ap.SetSelling(orderbook.Asset(txs.Transactions[i].SellType))

		// insert transaction into orderbook
		var pcs orderbook.PriceCompStats
		pcs.SellPrice = txs.Transactions[i].SellLimitPrice
		pcs.Txid = txs.Transactions[i].Txid
		pcs.CumulativeOfferedForSale = txs.Transactions[i].SellAmount
		pcs.CumulativeOfferedForSaleTimesPrice = txs.Transactions[i].SellAmount * txs.Transactions[i].SellLimitPrice

		if obval, ok := obm.MOrderbooks[ap]; ok { // if AssetPair already in OrderBookManager
			pos := sort.Search(len(obval.MPrecomputedTatonnementData), func(i int) bool { return obval.MPrecomputedTatonnementData[i].SellPrice >= pcs.SellPrice })
			newobval := make([]orderbook.PriceCompStats, len(obval.MPrecomputedTatonnementData)+1)
			at := copy(newobval, obval.MPrecomputedTatonnementData[:pos])
			newobval[pos] = pcs
			at++
			copy(newobval[at:], obval.MPrecomputedTatonnementData[pos:])

			for j := pos + 1; j < len(newobval); j++ {
				newobval[j].CumulativeOfferedForSale += pcs.CumulativeOfferedForSale
				newobval[j].CumulativeOfferedForSaleTimesPrice += pcs.CumulativeOfferedForSaleTimesPrice
			}

			newpcs := orderbook.Orderbook{MPrecomputedTatonnementData: newobval}
			obm.MOrderbooks[ap] = newpcs
		} else { // if AssetPair not already in OrderBookManager then add
			var ob orderbook.Orderbook
			ob.MPrecomputedTatonnementData = append(ob.MPrecomputedTatonnementData, pcs)
			obm.MOrderbooks[ap] = ob
		}
	}
	return obm
}

func main() {

	ms := new(tatonnement.TatonnementOracle)

	cp := tatonnement.ControlParams{MTaxRate: 0, MSmoothMult: 10, MMaxRounds: 3,
		MMStepUp: 0, MMStepDown: 0, MStepSizeRadix: 0,
		MStepRadix: 100000000, KStartingStepSize: 1, KPriceMin: 1, KPriceMax: 10000000, MRoundNumber: 1}
	prices := make(map[orderbook.Asset]float64)

	prices["ETH"] = 4500
	prices["USDT"] = 1

	jsonFile, err := os.Open("test_cases/txs.json")
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Successfully Opened txs.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var txs orderbook.Transactions
	json.Unmarshal(byteValue, &txs)

	obm := createOrderbookManager(&txs)

	ms.MOrderbookManager = *obm
	ms.ComputePrices(cp, prices)

}
