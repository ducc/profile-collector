package prom

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CPUNanosecondsSum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_nanoseconds_sum",
		Help: "Sum of cpu nanoseconds",
	}, []string{"function"})

	CPUNanosecondsMean = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_nanoseconds_mean",
		Help: "Mean of cpu nanoseconds",
	}, []string{"function"})

	CPUProfileSamplesSum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_profile_samples_sum",
		Help: "Sum of cpu profile samples",
	}, []string{"function"})

	CPUProfileSamplesMean = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "cpu_profile_samples_mean",
		Help: "Mean of cpu profile samples",
	}, []string{"function"})
)
