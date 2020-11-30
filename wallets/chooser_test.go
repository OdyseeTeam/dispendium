package wallets

import (
	"math/rand"
	"testing"
	"time"

	"github.com/lbryio/lbry.go/v2/lbrycrd"
)

func TestChooseWallet(t *testing.T) {
	rand.Seed(time.Now().Unix())
	AddWallet("a", &lbrycrd.Client{})
	AddWallet("b", &lbrycrd.Client{})
	AddWallet("c", &lbrycrd.Client{})

	for i := 0; i < 100; i++ {
		wallet := ChooseWallet()
		if wallet == nil {
			t.Error("no wallet!")
		}
		println("found wallet ", wallet.Name)
	}
}
