package jobs

import (
	"github.com/lbryio/dispendium/wallets"

	"github.com/jasonlvhit/gocron"
	"github.com/sirupsen/logrus"
)

var cronRunning chan bool
var scheduler *gocron.Scheduler

//MinLBCBalance is the minumum amount of LBC the LBRYcrd wallet can have before warnings are displayed. If slack integration
// is turned on then it will send the warning to slack as well.
var MinLBCBalance float64

//CreditBalanceLogger is a special logger used to to hook into a specific channel when the balance is below a certain
//threshold.
var CreditBalanceLogger *logrus.Logger

// Start starts the jobs that run in the background after initialization
func Start() {
	scheduler = gocron.NewScheduler()
	scheduler.Every(10).Minutes().From(gocron.NextTick()).Do(WalletBalanceCheck)

	cronRunning = scheduler.Start()
}

// Shutdown is used to shutdown the background jobs.
func Shutdown() {
	logrus.Debug("Shutting down cron jobs...")
	scheduler.Clear()
	close(cronRunning)
}

//WalletBalanceCheck checks the wallet balance and warns when its too low via warning
func WalletBalanceCheck() {
	balances, err := wallets.GetBalances()
	if err != nil {
		logrus.Error(err)
		return
	}
	addresses, err := wallets.GetAddresses()
	if err != nil {
		logrus.Error(err)
		return
	}
	addrMap := make(map[string]string)
	for _, addressList := range addresses {
		for _, addy := range addressList.Addresses {
			addrMap[addressList.Name] = addy
		}
	}
	for _, balance := range balances {
		if balance.LBC < MinLBCBalance {
			if CreditBalanceLogger != nil {
				CreditBalanceLogger.Warningf("ALERT: Dispendium balance(%.2f) is below the min balance of %.2f for %s - %s", balance.LBC, MinLBCBalance, balance.Name, addrMap[balance.Name])
			}
			logrus.Warningf("ALERT: Dispendium balance(%.2f) is below the min balance of %.2f for %s - %s", balance.LBC, MinLBCBalance, balance.Name, addrMap[balance.Name])
		}
	}
}
