package config

import (
	"github.com/lbryio/dispendium/wallets"
	"strconv"

	"github.com/lbryio/dispendium/jobs"

	"github.com/lbryio/dispendium/actions"
	"github.com/lbryio/dispendium/env"
	"github.com/lbryio/dispendium/util"

	"github.com/lbryio/lbry.go/v2/lbrycrd"

	"github.com/johntdyer/slackrus"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// InitializeConfiguration inits the base configuration of dispendium
func InitializeConfiguration() {
	conf, err := env.NewWithEnvVars()
	if err != nil {
		logrus.Panic(err)
	}
	if viper.GetBool("debugmode") {
		util.Debugging = true
		logrus.SetLevel(logrus.DebugLevel)
	}
	if viper.GetBool("tracemode") {
		util.Debugging = true
		logrus.SetLevel(logrus.TraceLevel)
	}
	util.AuthToken = conf.AuthToken
	initSlack(conf)
	initWallets(conf)
	SetLBCConfig(conf)
}

//SetLBCConfig sets the configuration for any environment variables related to the spending of LBC
func SetLBCConfig(config *env.Config) {
	var err error
	actions.MaxLBCPerHour, err = strconv.ParseFloat(config.MaxLBCPerHour, 64)
	if err != nil {
		logrus.Panic(err)
	}
	actions.MaxLBCPayment, err = strconv.ParseFloat(config.MaxLBCPayment, 64)
	if err != nil {
		logrus.Panic(err)
	}
	jobs.MinLBCBalance, err = strconv.ParseFloat(config.MinBalance, 64)
	if err != nil {
		logrus.Panic(err)
	}
}

// initSlack initializes the slack connection and posts info level or greater to the set channel.
func initSlack(config *env.Config) {
	slackURL := config.SlackHookURL
	slackChannel := config.SlackChannel
	if slackURL != "" && slackChannel != "" {
		logrus.AddHook(&slackrus.SlackrusHook{
			HookURL:        slackURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
			Channel:        slackChannel,
			IconEmoji:      ":money_mouth_face:",
			Username:       "Dispendium",
		})

		jobs.CreditBalanceLogger = logrus.New()
		jobs.CreditBalanceLogger.AddHook(&slackrus.SlackrusHook{
			HookURL:        slackURL,
			AcceptedLevels: slackrus.LevelThreshold(logrus.InfoLevel),
			Channel:        "credit-alerts",
			IconEmoji:      ":money_mouth_face:",
			Username:       "wallet-watcher",
		})
	}
}

func initWallets(conf *env.Config) {
	chainParams, ok := lbrycrd.ChainParamsMap[conf.BlockchainName]
	if !ok {
		logrus.Panicf("block chain name %s is not recognized", conf.BlockchainName)
	}
	wallets.SetChainParams(&chainParams)
	instances := viper.GetStringMapString("lbrycrd")
	for name, url := range instances {
		lbrycrdClient, err := lbrycrd.New(url, &chainParams)
		if err != nil {
			panic(err)
		}
		_, err = lbrycrdClient.GetBalance("*")
		if err != nil {
			logrus.Panicf("Error connecting to lbrycrd: %+v", err)
		}
		wallets.AddWallet(name, lbrycrdClient)
	}
}
