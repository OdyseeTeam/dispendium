package util

//Debugging lets the application know if it is in debugging mode.
var Debugging bool

//LBC is a utility function to convert satoshi amounts to LBC amounts.
func LBC(satoshiAmount uint64) float64 {
	return float64(satoshiAmount) / 100000000
}
