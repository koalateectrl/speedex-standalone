package demandutils

import (
	"math"
	"strconv"

	"github.com/sandymule/speedex-standalone/pkg/assets"
)

type SupplyDemandPair struct {
	Supply float64
	Demand float64
}

type SupplyDemand struct {
	MSupplyDemand map[assets.Asset]*SupplyDemandPair
}

type TatonnementObjectiveFunction struct {
	Value float64
}

func (sd *SupplyDemand) AddSupplyDemandPair(tradingPair assets.AssetPair, amount float64) {
	sd.AddSupplyDemand(tradingPair.Selling(), tradingPair.Buying(), amount)
}

func (sd *SupplyDemand) AddSupplyDemand(sell assets.Asset, buy assets.Asset, amount float64) {
	if _, found := sd.MSupplyDemand[sell]; found {
		sd.MSupplyDemand[sell].Supply += amount
	} else {
		sd.MSupplyDemand[sell] = &SupplyDemandPair{amount, 0}
	}

	if _, found := sd.MSupplyDemand[buy]; found {
		sd.MSupplyDemand[buy].Demand += amount
	} else {
		sd.MSupplyDemand[buy] = &SupplyDemandPair{0, amount}
	}
}

func (sd *SupplyDemand) GetDelta(asset assets.Asset) float64 {
	return sd.MSupplyDemand[asset].Demand - sd.MSupplyDemand[asset].Supply
}

func (sd *SupplyDemand) GetObjective() *TatonnementObjectiveFunction {
	return NewTatonnementObjectiveFunction(sd)
}

func NewTatonnementObjectiveFunction(supplyDemand *SupplyDemand) *TatonnementObjectiveFunction {
	tof := new(TatonnementObjectiveFunction)
	for _, sd := range supplyDemand.MSupplyDemand {
		tof.Value += math.Pow(sd.Demand-sd.Supply, 2)
	}

	return tof
}

func (sd *SupplyDemand) String() string {
	retStr := "{"
	for key, val := range sd.MSupplyDemand {
		retStr += string(key) + ": ("
		retStr += strconv.FormatFloat(val.Supply, 'f', -1, 64) + "," + strconv.FormatFloat(val.Demand, 'f', -1, 64) + "); "
	}
	retStr += "}"
	return retStr
}
