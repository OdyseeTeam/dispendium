package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (

	//SendAPI tracks metrics related to Send API of dispendium. There are intended to be two types, duration and amount. They
	//track the spending output and how long it takes for the api to complete. This is used with the LBRY monitors
	SendAPI = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "dispendium",
		Subsystem: "send",
		Name:      "api",
		Help:      "Tracks the send amounts and send durations for the api",
	}, []string{"type"})
)
