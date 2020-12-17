package actions

import (
	"net/http"
	"time"

	"github.com/lbryio/dispendium/dispendiumapi"
	"github.com/lbryio/dispendium/internal/metrics"
	"github.com/lbryio/dispendium/util"
	"github.com/lbryio/dispendium/wallets"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/lbrycrd"
	v "github.com/lbryio/ozzo-validation"
	"github.com/lbryio/ozzo-validation/is"

	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"
)

// MaxLBCPerHour is the maximum amount of LBC that can be sent in an hour.
var MaxLBCPerHour float64

// MaxLBCPayment is the maximum a payment can be
var MaxLBCPayment float64

// Root Handler is the default handler
func Root(r *http.Request) api.Response {
	if r.Method != http.MethodPost {
		return api.Response{Error: errors.Err("invalid method"), Status: http.StatusMethodNotAllowed}
	}
	if r.URL.Path == "/" {
		return api.Response{Data: "Dispendium, dispenser of rewards"}
	}
	return api.Response{Status: http.StatusNotFound, Error: errors.Err("404 Not Found")}
}

// Test Handler can be used for testing and triggering
func Test(r *http.Request) api.Response {
	return api.Response{Data: "ok"}
}

// Send Handler sends LBC via LBRYcrd to a wallet addressed passed in.
func Send(r *http.Request) api.Response {
	if r.Method != http.MethodPost && !util.Debugging {
		return api.Response{Error: errors.Err("invalid method"), Status: http.StatusMethodNotAllowed}
	}
	start := time.Now()
	defer metrics.SendAPI.WithLabelValues("duration").Observe(time.Since(start).Seconds())
	params := struct {
		AuthToken     string
		WalletAddress string
		SatoshiAmount int64
	}{}
	err := api.FormValues(r, &params, []*v.FieldRules{
		v.Field(&params.AuthToken, v.Required),
		v.Field(&params.WalletAddress, v.Required, is.Alphanumeric),
		v.Field(&params.SatoshiAmount, v.Required, v.Min(0)),
	})
	if err != nil {
		return api.Response{Error: err, Status: http.StatusBadRequest}
	}

	if params.AuthToken != util.AuthToken {
		logrus.Warningf("Login with incorrect token %s", params.AuthToken)
		return api.Response{Error: errors.Err("not authorized"), Status: http.StatusUnauthorized}
	}

	if util.LBC(uint64(params.SatoshiAmount)) > MaxLBCPayment {
		return api.Response{Error: errors.Err("(258200) Sending disabled. Amount cannot be sent. Please reach out to support@lbry.com."), Status: http.StatusBadRequest}
	}

	decodedAddress, err := lbrycrd.DecodeAddress(params.WalletAddress, wallets.GetCainParams())
	if err != nil {
		return api.Response{Error: errors.Err("could not decode wallet address, please check network and chain: ", err)}
	}
	amount := btcutil.Amount(params.SatoshiAmount)
	wallet, err := wallets.ChooseWallet()
	if err == nil {
		return api.Response{Error: err}
	}
	txHash, err := wallet.SendToAddress(decodedAddress, amount)
	if err != nil {
		wallets.RemoveWallet(wallet)
		return api.Response{Error: errors.Err(err)}
	}

	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}
	metrics.SendAPI.WithLabelValues("amount").Observe(util.LBC(uint64(params.SatoshiAmount)))
	logrus.Debugf("Sending %.2f LBC to %s", util.LBC(uint64(params.SatoshiAmount)), params.WalletAddress)
	return api.Response{Data: dispendiumapi.SendResult{
		LBCAmount:     util.LBC(uint64(params.SatoshiAmount)),
		SatoshiAmount: uint64(params.SatoshiAmount),
		TxHash:        txHash.String(),
	}}
}

// Balance Handler returns lbrycrd wallet available balances
func Balance(r *http.Request) api.Response {
	if r.Method != http.MethodPost && !util.Debugging {
		return api.Response{Error: errors.Err("invalid method"), Status: http.StatusMethodNotAllowed}
	}
	params := struct {
		AuthToken string
	}{}
	err := api.FormValues(r, &params, []*v.FieldRules{
		v.Field(&params.AuthToken, v.Required),
	})
	if err != nil {
		return api.Response{Error: err, Status: http.StatusBadRequest}
	}

	if params.AuthToken != util.AuthToken {
		logrus.Warningf("Login with incorrect token %s", params.AuthToken)
		return api.Response{Error: errors.Err("not authorized"), Status: http.StatusUnauthorized}
	}
	balances, err := wallets.GetBalances()
	if err != nil {
		return api.Response{Error: err}
	}

	return api.Response{Data: balances}
}

// Addresses Handler returns lbrycrd wallet available balances
func Addresses(r *http.Request) api.Response {
	if r.Method != http.MethodPost && !util.Debugging {
		return api.Response{Error: errors.Err("invalid method"), Status: http.StatusMethodNotAllowed}
	}
	params := struct {
		AuthToken string
	}{}
	err := api.FormValues(r, &params, []*v.FieldRules{
		v.Field(&params.AuthToken, v.Required),
	})
	if err != nil {
		return api.Response{Error: err, Status: http.StatusBadRequest}
	}

	if params.AuthToken != util.AuthToken {
		logrus.Warningf("Login with incorrect token %s", params.AuthToken)
		return api.Response{Error: errors.Err("not authorized"), Status: http.StatusUnauthorized}
	}
	addresses, err := wallets.GetAddresses()
	if err != nil {
		return api.Response{Error: err}
	}

	return api.Response{Data: addresses}
}
