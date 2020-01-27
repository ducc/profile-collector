package main

import (
	"github.com/ducc/profile-collector/store/aggregator/prom"
	pprof "github.com/google/pprof/profile"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// ctx := context.Background()

	fileName := "/Users/joeburnard/Desktop/collector_cpu_1.prof"

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	profile, err := pprof.ParseData(data)
	if err != nil {
		panic(err)
	}

	for _, valueType := range profile.SampleType {
		logrus.Info(valueType.Unit, " ", valueType.Type)
	}

	// sampleType0 := profile.SampleType[0].Type
	// sampleUnit0 := profile.SampleType[0].Unit
	//
	// sampleType1 := profile.SampleType[1].Type
	// sampleUnit1 := profile.SampleType[1].Unit
	//
	// for i, sample := range profile.Sample {
	// 	logrus.Warnf("")
	// 	logrus.Warnf("== SAMPLE %d ==", i+1)
	// 	logrus.Warnf("%s %s %d", profile.PeriodType.Type, profile.PeriodType.Unit, profile.Period)
	// 	logrus.Warn(sample.Value)
	// 	for _, location := range sample.Location {
	// 		for _, line := range location.Line {
	// 			logrus.Infof("%s:%d - %s (%s): %d - %s (%s): %d", line.Function.Name, line.Line, sampleType0, sampleUnit0, sample.Value[0], sampleType1, sampleUnit1, sample.Value[1])
	// 		}
	// 	}
	// }

	p, _ := prom.New()
	p.PerformAggregateOfSamples(profile, prom.MeanAggregator)

	// profile := &protos.Profile{}
	//
	// if err := proto.Unmarshal(decompressed, profile); err != nil {
	// 	panic(err)
	// }

	logrus.Info("yes we have the profile")
}

type IProfile struct {
	Samples []Sample
}

type Sample struct {
	Location struct {
	}
}
