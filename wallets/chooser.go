package wallets

import (
	"math/rand"
)

// ChooseWallet is currently a round robin
func ChooseWallet() *Wallet {
	return &loadedWallets[int(rand.Float64()*100)%len(loadedWallets)]
}
