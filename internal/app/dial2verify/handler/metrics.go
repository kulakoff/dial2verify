package handler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	checkPhoneRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dial2verify_check_phone_requests_total",
		Help: "Total number of /api/checkPhone requests received",
	})

	checkPhoneErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dial2verify_check_phone_errors_total",
		Help: "Total number of /api/checkPhone requests that resulted in an error",
	})

	checkPhoneInvalidTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dial2verify_check_phone_invalid_total",
		Help: "Total number of /api/checkPhone requests with invalid phone format",
	})

	checkPhoneDurationSeconds = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "dial2verify_check_phone_duration_seconds",
		Help:    "Duration of /api/checkPhone handler in seconds",
		Buckets: prometheus.DefBuckets,
	})

	checkPhoneFoundTrueTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dial2verify_check_phone_found_true_total",
		Help: "Total number of /api/checkPhone responses where found=true",
	})

	checkPhoneFoundFalseTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "dial2verify_check_phone_found_false_total",
		Help: "Total number of /api/checkPhone responses where found=false",
	})
)
