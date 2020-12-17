package wallets

import (
	"math/rand"

	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// ChooseWallet is currently a round robin
func ChooseWallet() (*Wallet, error) {
	if len(loadedWallets) == 0 {
		return nil, errors.Err("there are no loaded wallets for dispendium!")
	}
	return &loadedWallets[int(rand.Float64()*100)%len(loadedWallets)], nil
}
