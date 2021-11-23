package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sandymule/speedex-standalone/orderbook"
	"github.com/sandymule/speedex-standalone/tatonnement"
)

func main() {
	ap := new(orderbook.AssetPair)
	ap.SetBuying("ETH")
	ap.SetSelling("USDC")

	ms := new(tatonnement.TatonnementOracle)
	ms.SetStr1("Sam")
	ms.SetStr2("Wong")

	jsonFile, err := os.Open("test_cases/txs.json")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened txs.json")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var txs orderbook.Transactions
	json.Unmarshal(byteValue, &txs)

	for i := 0; i < len(txs.Transactions); i++ {
		fmt.Println(txs.Transactions[i].Txid)
		fmt.Println(txs.Transactions[i].BuyAsset.Type)
		fmt.Println(txs.Transactions[i].BuyAsset.Amount)
		fmt.Println(txs.Transactions[i].BuyAsset.LimitPrice)

		fmt.Println(txs.Transactions[i].SellAsset.Type)
		fmt.Println(txs.Transactions[i].SellAsset.Amount)
		fmt.Println(txs.Transactions[i].SellAsset.LimitPrice)

	}

}
