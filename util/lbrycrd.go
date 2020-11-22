package util

import (
	"github.com/btcsuite/btcd/chaincfg"

	"github.com/lbryio/lbry.go/v2/lbrycrd"
)

//LbrycrdClient client for lbrycrd to be used in the app
var LbrycrdClient *lbrycrd.Client

//ChainParams chain parameters used in the application
var ChainParams *chaincfg.Params
