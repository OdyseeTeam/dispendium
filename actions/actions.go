package actions

import (
	"net/http"
	"time"

	"github.com/lbryio/dispendium/internal/metrics"
	"github.com/lbryio/dispendium/util"

	"github.com/lbryio/lbry.go/v2/extras/api"
	"github.com/lbryio/lbry.go/v2/extras/errors"
	"github.com/lbryio/lbry.go/v2/extras/null"
	"github.com/lbryio/lbry.go/v2/lbrycrd"
	v "github.com/lbryio/ozzo-validation"
	"github.com/lbryio/ozzo-validation/is"

	"github.com/btcsuite/btcutil"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"
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
	if r.Method != http.MethodPost {
		return api.Response{Error: errors.Err("invalid method"), Status: http.StatusMethodNotAllowed}
	}
	start := time.Now()
	defer metrics.SendAPI.WithLabelValues("duration").Observe(time.Since(start).Seconds())
	params := struct {
		AuthToken     string
		WalletAddress string
		SatoshiAmount int64
		CustomerID    *string
	}{}
	err := api.FormValues(r, &params, []*v.FieldRules{
		v.Field(&params.AuthToken, v.Required),
		v.Field(&params.WalletAddress, v.Required, is.Alphanumeric),
		v.Field(&params.SatoshiAmount, v.Required, v.Min(0)),
		v.Field(&params.CustomerID, is.PrintableASCII),
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

	if err := withinLBCLimits(); err != nil {
		return api.Response{Error: errors.Err(err)}
	}
	decodedAddress, err := lbrycrd.DecodeAddress(params.WalletAddress, util.ChainParams)
	if err != nil {
		return api.Response{Error: errors.Err("could not decode wallet address, please check network and chain: ", err)}
	}
	amount := btcutil.Amount(params.SatoshiAmount)
	txHash, err := util.LbrycrdClient.SendToAddress(decodedAddress, amount)
	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}

	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}
	metrics.SendAPI.WithLabelValues("amount").Observe(util.LBC(uint64(params.SatoshiAmount)))
	logrus.Debugf("Sending %.2f LBC to %s", util.LBC(uint64(params.SatoshiAmount)), params.WalletAddress)
	return api.Response{Data: struct {
		LBCAmount float64 `json:"lbc_amount"`
		TxHash    string  `json:"tx_id"`
	}{
		util.LBC(uint64(params.SatoshiAmount)),
		txHash.String(),
	}}
}

func withinLBCLimits() error {
	//payments, err := model.Payments(qmhelper.Where(model.PaymentColumns.CreatedAt, qmhelper.GTE,time.Now().Add(-1 *time.Hour))).AllG()
	result := boil.GetDB().QueryRow(`SELECT SUM(payment.satoshi_amount) FROM payment WHERE payment.created_at >= ?`, time.Now().Add(-1*time.Hour))
	var satoshiPaid null.Uint64
	err := result.Scan(&satoshiPaid)
	if err != nil {
		return errors.Err(err)
	}

	if util.LBC(satoshiPaid.Uint64) >= MaxLBCPerHour {
		return api.StatusError{Err: errors.Err("(244600) Sending disabled. Amount cannot be sent. Please reach out to support@lbry.com."), Status: http.StatusBadRequest}
	} else if util.LBC(satoshiPaid.Uint64)/MaxLBCPerHour > 0.80 {
		logrus.Warnf("Send API is within 20%% its rate limit of %f LBC per hour and will automatically disable", MaxLBCPerHour)
	}

	return nil
}

type balance struct {
	LBC     float64 `json:"lbc"`
	Satoshi uint64  `json:"satoshi"`
}

// Balance Handler returns lbrycrd available balance
func Balance(r *http.Request) api.Response {
	if r.Method != http.MethodPost {
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

	available, err := util.LbrycrdClient.GetBalanceMinConf("*", 1)
	if err != nil {
		return api.Response{Error: errors.Err(err)}
	}

	balance := &balance{
		LBC:     available.ToBTC(),
		Satoshi: uint64(available),
	}

	return api.Response{Data: balance}
}
