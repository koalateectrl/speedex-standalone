package tatonnement

type ControlParams struct {
	MTaxRate, MSmoothMult                uint8
	MMaxRounds                           uint32
	MMStepUp, MMStepDown, MStepSizeRadix uint8
	MStepRadix                           uint8

	MRoundNumber uint32
}

func (cp *ControlParams) IncrementRound() {
	cp.MRoundNumber++
}

func (cp *ControlParams) Done() bool {
	return cp.MRoundNumber >= cp.MMaxRounds
}

func (cp *ControlParams) SetTrialPrice(curPrice uint64, demand int64, stepSize uint64) {
	// set price for one asset
}

func (cp *ControlParams) SetTrialPrices() {
	//set prices for all assets
}
