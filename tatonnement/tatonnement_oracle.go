package tatonnement

type TatonnementOracle struct {
	str1   string
	str2   string
	Params ControlParams
}

func (to *TatonnementOracle) SetStr1(newstr1 string) { // setter
	to.str1 = newstr1
}

func (to *TatonnementOracle) SetStr2(newstr2 string) { // setter
	to.str2 = newstr2
}

func (to *TatonnementOracle) String() string {
	return "(" + to.str1 + " / " + to.str2 + ")"
}
