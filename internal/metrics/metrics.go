package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (

	//SendAmount tracks metrics related to Send API of dispendium. It tracks the spending output.
	SendAmount = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "dispendium",
		Subsystem: "sending",
		Name:      "amount",
		Help:      "Tracks the send amounts for the api",
	}, []string{"instance"})

	//sendDur tracks metrics related to Send API of Dispendium. It tracks how long it takes for the sending to complete.
	sendDur = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "dispendium",
		Subsystem: "sending",
		Name:      "duration",
		Help:      "Tracks the send durations for the api",
	}, []string{"instance"})

	//apiDuration It tracks how long it takes for the api's of Dispendium to complete.
	apiDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "dispendium",
		Subsystem: "api",
		Name:      "duration",
		Help:      "Tracks the durations for an api",
	}, []string{"api"})

	// Balance tracks the balance of each wallet whenever we call get balance.
	Balance = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "dispendium",
		Subsystem: "lbrycrd",
		Name:      "balance",
		Help:      "Tracks the send durations for the api",
	}, []string{"instance"})
)

// SendDuration logs the time it takes to send with the respective metric
func SendDuration(start time.Time, instance string) {
	duration := time.Since(start).Seconds()
	sendDur.WithLabelValues(instance).Observe(duration)
}

// APIDuration logs the time it takes to complete the api with the respective metric
func APIDuration(start time.Time, api string) {
	duration := time.Since(start).Seconds()
	apiDuration.WithLabelValues(api).Observe(duration)
}
