package wallets

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/lbryio/lbry.go/v2/lbrycrd"
)

// Wallet an instance of a wallet that can be used by dispendium, equivalent to a lbrycrd instance running
type Wallet struct {
	*lbrycrd.Client
	Name string
}

//LbrycrdClient client for lbrycrd to be used in the app
var loadedWallets []Wallet

//chainParams chain parameters used in the application
var chainParams *chaincfg.Params

// SetChainParams sets the chain paramters used by dispendium for validating addresses and making calls to lbrycrd instances
func SetChainParams(params *chaincfg.Params) {
	chainParams = params
}

// GetCainParams retreives the chain params set on initialization
func GetCainParams() *chaincfg.Params {
	return chainParams
}

// AddWallet adds a wallet to the loaded wallets that are to be used by dispendium
func AddWallet(name string, client *lbrycrd.Client) {
	loadedWallets = append(loadedWallets, Wallet{
		Client: client,
		Name:   name,
	})
}

// Balance balance of a wallet instance used by dispendium
type Balance struct {
	Name    string  `json:"name"`
	LBC     float64 `json:"lbc"`
	Satoshi uint64  `json:"satoshi"`
}

// GetBalances retrieves the balances for all wallet instances used by dispendium
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

// WalletAccount accounts used by a wallet lbrycrd instance. It holds the available addresses for sending to dispendium
type WalletAccount struct {
	Name      string   `json:"name"`
	Addresses []string `json:"addresses"`
}

// GetAddresses returns the set of addresses in use for each wallet lbrycrd instance loaded into dispendium
func GetAddresses() ([]WalletAccount, error) {
	var addresses []WalletAccount
	for _, wallet := range loadedWallets {
		address := WalletAccount{Name: wallet.Name}
		results, err := wallet.ListReceivedByAddressIncludeEmpty(0, true)
		if err != nil {
			return nil, errors.Err(err)
		}
		for _, account := range results {
			address.Addresses = append(address.Addresses, account.Address)
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}
