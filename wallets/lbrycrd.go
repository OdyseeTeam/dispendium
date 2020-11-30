package wallets

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/lbry.go/v2/lbrycrd"
)

type Wallet struct {
	*lbrycrd.Client
	Name string
}

//LbrycrdClient client for lbrycrd to be used in the app
var loadedWallets []Wallet

//chainParams chain parameters used in the application
var chainParams *chaincfg.Params

func SetChainParams(params *chaincfg.Params) {
	chainParams = params
}

func GetCainParams() *chaincfg.Params {
	return chainParams
}

func AddWallet(name string, client *lbrycrd.Client) {
	loadedWallets = append(loadedWallets, Wallet{
		Client: client,
		Name:   name,
	})
}

type Balance struct {
	Name    string  `json:"name"`
	LBC     float64 `json:"lbc"`
	Satoshi uint64  `json:"satoshi"`
}

func GetBalances() ([]Balance, error) {
	var balances []Balance
	for _, wallet := range loadedWallets {
		available, err := wallet.GetBalanceMinConf("*", 1)
		if err != nil {
			return nil, errors.Err(err)
		}

		balances = append(balances, Balance{
			Name:    wallet.Name,
			LBC:     available.ToBTC(),
			Satoshi: uint64(available),
		})
	}
	return balances, nil
}
