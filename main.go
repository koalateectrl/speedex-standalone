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
		ap.SetBuying(orderbook.Asset(txs.Transactions[i].BuyAsset.Type))
		ap.SetSelling(orderbook.Asset(txs.Transactions[i].SellAsset.Type))

		// insert transaction into orderbook
		var pcs orderbook.PriceCompStats
		pcs.SellPrice = txs.Transactions[i].SellAsset.LimitPrice
		pcs.BuyPrice = txs.Transactions[i].BuyAsset.LimitPrice
		pcs.Txid = txs.Transactions[i].Txid
		pcs.CumulativeOfferedForSale = txs.Transactions[i].SellAsset.Amount

		if obval, ok := obm.MOrderbooks[ap]; ok {
			pos := sort.Search(len(obval.MPrecomputedTatonnementData), func(i int) bool { return obval.MPrecomputedTatonnementData[i].SellPrice >= pcs.SellPrice })
			newobval := make([]orderbook.PriceCompStats, len(obval.MPrecomputedTatonnementData)+1)
			at := copy(newobval, obval.MPrecomputedTatonnementData[:pos])
			newobval[pos] = pcs
			at++
			copy(newobval[at:], obval.MPrecomputedTatonnementData[pos:])

			for j := pos + 1; j < len(newobval); j++ {
				newobval[j].CumulativeOfferedForSale += txs.Transactions[i].SellAsset.Amount
			}

			newpcs := orderbook.Orderbook{MPrecomputedTatonnementData: newobval}
			obm.MOrderbooks[ap] = newpcs
		} else {
			var ob orderbook.Orderbook
			ob.MPrecomputedTatonnementData = append(ob.MPrecomputedTatonnementData, pcs)
			obm.MOrderbooks[ap] = ob
		}
	}
	return obm
}

func main() {

	ms := new(tatonnement.TatonnementOracle)
	fmt.Println(ms)

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
	fmt.Println(obm.MOrderbooks)

}
