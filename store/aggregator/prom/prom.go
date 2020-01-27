package prom

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/golang/protobuf/proto"
	pprof "github.com/google/pprof/profile"
	"github.com/pkg/errors"
)

type aggregator struct {
}

func New() (*aggregator, error) {
	return &aggregator{}, nil
}

func (a *aggregator) Add(ctx context.Context, storedProfile *protos.StoredProfile) error {
	data, err := proto.Marshal(storedProfile.Profile)
	if err != nil {
		return err
	}

	prof, err := pprof.ParseUncompressed(data)
	if err != nil {
		return err
	}

	// gross
	switch storedProfile.Metadata.ProfileType {
	case protos.ProfileMetadata_TYPE_CPU:
		aggregated := a.AggregateCPUSamples(prof)
		a.UpdateCPUMetrics(aggregated)
	default:
		return errors.Errorf("unsupported profile type: %s", storedProfile.Metadata.ProfileType)
	}

	return nil
}

type AggregatedCPUSample struct {
	samples     float64
	samplesSum  float64
	samplesMean float64
	timeSum     float64
	timeMean    float64
}

type CPUAggregator map[string]*AggregatedCPUSample

func (c CPUAggregator) Put(functionName string, samples float64, time float64) {
	agg, ok := c[functionName]
	if !ok {
		c[functionName] = &AggregatedCPUSample{}
	}

	agg.samples++

	agg.samplesMean = (agg.samplesSum + samples) / agg.samples
	agg.samplesSum += samples

	agg.timeMean = (agg.timeSum + time) / agg.samples
	agg.timeSum += time
}

func (a *aggregator) AggregateCPUSamples(prof *pprof.Profile) map[string]*AggregatedCPUSample {
	agg := CPUAggregator(make(map[string]*AggregatedCPUSample))

	for _, sample := range prof.Sample {
		for _, location := range sample.Location {
			for _, line := range location.Line {
				agg.Put(line.Function.Name, float64(sample.Value[0]), float64(sample.Value[1]))
			}
		}
	}

	return agg
}

func (a *aggregator) UpdateCPUMetrics(aggregated map[string]*AggregatedCPUSample) {
	for functionName, sample := range aggregated {
		CPUNanosecondsSum.WithLabelValues(functionName).Set(sample.timeSum)
		CPUNanosecondsMean.WithLabelValues(functionName).Set(sample.timeMean)

		CPUProfileSamplesSum.WithLabelValues(functionName).Set(sample.samplesSum)
		CPUProfileSamplesMean.WithLabelValues(functionName).Set(sample.samplesMean)
	}
}
