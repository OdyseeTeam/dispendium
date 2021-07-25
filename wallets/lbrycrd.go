package wallets

import (
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcutil"
	"github.com/lbryio/dispendium/dispendiumapi"
	"github.com/lbryio/dispendium/internal/metrics"
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

// GetBalances retrieves the balances for all wallet instances used by dispendium
func GetBalances() ([]dispendiumapi.BalanceResult, error) {
	var balances []dispendiumapi.BalanceResult
	for _, wallet := range loadedWallets {
		available, err := wallet.GetBalanceMinConf("*", 0)
		if err != nil {
			return nil, errors.Err(err)
		}
		metrics.Balance.WithLabelValues(wallet.Name).Set(available.ToBTC())
		balances = append(balances, dispendiumapi.BalanceResult{
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

// RemoveWallet removes a wallet from rotation!
func RemoveWallet(w *Wallet) {
	var newSet []Wallet
	for _, wallet := range loadedWallets {
		if wallet.Name != w.Name {
			newSet = append(newSet, wallet)
		}
	}
	loadedWallets = newSet
}

func (c *Wallet) SendToAddress(address btcutil.Address, amount btcutil.Amount) (*chainhash.Hash, error) {
	defer metrics.SendDuration(time.Now(), c.Name)
	defer metrics.SendAmount.WithLabelValues(c.Name).Observe(amount.ToBTC())
	return c.Client.SendToAddress(address, amount)
}
