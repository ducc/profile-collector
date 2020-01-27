package profiles

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/ducc/profile-collector/store/profiles/gcs"
)

type Profiles interface {
	Add(ctx context.Context, profile *protos.StoredProfile) error
}

func NewGCS(ctx context.Context, bucket string) (Profiles, error) {
	return gcs.New(ctx, bucket, gcs.TimeFormatter)
}

type MockProfiles struct {
}

func (m *MockProfiles) Add(ctx context.Context, profile *protos.StoredProfile) error {
	return nil
}
