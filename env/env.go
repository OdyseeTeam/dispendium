package env

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"

	e "github.com/caarlos0/env"
)

// Config holds the environment configuration used by lighthouse.
type Config struct {
	LbrycrdURL     string `env:"LBRYCRD_CONNECT" envDefault:""`
	SlackHookURL   string `env:"SLACKHOOKURL"`
	SlackChannel   string `env:"SLACKCHANNEL"`
	AuthToken      string `env:"AUTH_TOKEN" envDefault:"MyTokeN"`
	BlockchainName string `env:"BLOCKCHAIN_NAME" envDefault:"lbrycrd_main"`
	MaxLBCPerHour  string `env:"MAX_LBC_PER_HR" envDefault:"100000"`
	MaxLBCPayment  string `env:"MAX_LBC_PAYMENT" envDefault:"50000"`
	MinBalance     string `env:"MIN_BALANCE" envDefault:"5000"`
}

// NewWithEnvVars creates an Config from environment variables
func NewWithEnvVars() (*Config, error) {
	cfg := &Config{}
	err := e.Parse(cfg)
	if err != nil {
		return nil, errors.Err(err)
	}

	return cfg, nil
}
