package dispendiumapi

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lbryio/lbry.go/v2/extras/errors"

	"github.com/btcsuite/btcutil"
)

type Client struct {
	*http.Client
	url   string
	token string
}

func NewClient(url, token string) *Client {
	return &Client{
		Client: &http.Client{},
		url:    url,
		token:  token,
	}
}

type SendFundsArgs struct {
	WalletAddress string
	Amount        btcutil.Amount
}

type SendResult struct {
	LBCAmount     float64 `json:"lbc_amount"`
	SatoshiAmount uint64  `json:"satoshi_amount"`
	TxHash        string  `json:"tx_id"`
}

type SendFundsResponse struct {
	*http.Response
	Success bool        `json:"success"`
	Error   *string     `json:"error"`
	Data    *SendResult `json:"data"`
	Trace   []string    `json:"_trace,omitempty"`
}

func (c Client) SendFunds(args SendFundsArgs) (*SendFundsResponse, error) {
	formData := url.Values{}
	formData.Add("auth_token", c.token)
	formData.Add("wallet_address", args.WalletAddress)
	formData.Add("satoshi_amount", strconv.Itoa(int(args.Amount.ToUnit(btcutil.AmountSatoshi))))
	resp, err := c.PostForm(c.url+"/send", formData)
	if err != nil {
		return nil, errors.Err(err)
	}
	defer resp.Body.Close()
	sfresp := &SendFundsResponse{}
	err = json.NewDecoder(resp.Body).Decode(sfresp)
	if err != nil {
		return nil, errors.Err(err)
	}
	return sfresp, nil
}
