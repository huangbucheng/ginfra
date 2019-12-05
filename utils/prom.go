package utils

import (
	"strconv"
	"time"

	mw "ginfra/middleware"

	"github.com/prometheus/client_golang/prometheus"
)

// fetch metrics and consume
// reference: https://github.com/prometheus/prom2json/blob/master/cmd/prom2json/main.go

func init() {
	mw.GRegistry.Register(SearchCount)
	mw.GRegistry.Register(SearchDuration)
	mw.GRegistry.Register(QueryError)
}

var SearchCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cloudsearch_apiv3_search",
		Help: "apiv3 search request count",
	},
	[]string{"appid"},
)

var SearchDuration = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Name: "cloudsearch_apiv3_search_duration",
		Help: "apiv3 search request duration",
	},
	[]string{"appid"},
)

var QueryError = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "cloudsearch_apiv3_error",
		Help: "apiv3 request error",
	},
	[]string{"appid", "action", "error"},
)

func SearchCountInc(appid uint64) {
	SearchCount.With(prometheus.Labels{
		"appid": strconv.FormatUint(appid, 10),
	}).Inc()
}

func SearchDurationObserve(appid uint64, begin time.Time) {
	SearchDuration.With(prometheus.Labels{
		"appid": strconv.FormatUint(appid, 10),
	}).Observe(float64(time.Since(begin)) / float64(time.Second))
}

func QueryErrorInc(appid uint64, action string, err error) {
	QueryError.With(prometheus.Labels{
		"appid":  strconv.FormatUint(appid, 10),
		"action": action,
		"error":  err.Error(),
	}).Inc()
}
