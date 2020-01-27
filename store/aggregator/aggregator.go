package aggregator

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/ducc/profile-collector/store/aggregator/prom"
)

type Aggregator interface {
	Add(ctx context.Context, profile *protos.StoredProfile) error
}

func NewPrometheus() (Aggregator, error) {
	return prom.New()
}
