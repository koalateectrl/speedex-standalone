package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/sandymule/speedex-standalone/pkg/assets"
	"github.com/sandymule/speedex-standalone/pkg/orderbook"
	"github.com/sandymule/speedex-standalone/pkg/tatonnement"
)

func createOrderbookManager(txs *orderbook.Transactions) *orderbook.OrderbookManager {
	obm := &orderbook.OrderbookManager{MOrderbooks: make(map[assets.AssetPair]orderbook.Orderbook)}
	for i := 0; i < len(txs.Transactions); i++ {
		var ap assets.AssetPair
		ap.SetBuying(assets.Asset(txs.Transactions[i].BuyType))
		ap.SetSelling(assets.Asset(txs.Transactions[i].SellType))

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
			// start with 0 price entry
			ob.MPrecomputedTatonnementData = append(ob.MPrecomputedTatonnementData, orderbook.PriceCompStats{SellPrice: 0, CumulativeOfferedForSale: 0, CumulativeOfferedForSaleTimesPrice: 0})
			ob.MPrecomputedTatonnementData = append(ob.MPrecomputedTatonnementData, pcs)
			obm.MOrderbooks[ap] = ob
		}
	}
	return obm
}

func main() {

	ms := new(tatonnement.TatonnementOracle)

	cp := tatonnement.TatonnementControlParams{MSmoothMult: 5, MMaxRounds: 50,
		MStepUp: 40, MStepDown: 25, MStepSizeRadix: 5, MStepRadix: 30}
	prices := make(map[assets.Asset]float64)

	prices["ETH"] = 4500
	prices["USDT"] = 1

	jsonFile, err := os.Open("../test_cases/txs.json")
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
	ms.ComputePrices(cp, prices, 5)

}
