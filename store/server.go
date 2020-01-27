package store

import (
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/ducc/profile-collector/store/aggregator"
	"github.com/ducc/profile-collector/store/profiles"
)

type server struct {
	profilesStore profiles.Profiles
	aggStore      aggregator.Aggregator
}

func New(profilesStore profiles.Profiles, aggStore aggregator.Aggregator) (protos.StoreServer, error) {
	return &server{
		profilesStore: profilesStore,
		aggStore:      aggStore,
	}, nil
}

func (s *server) ListProfiles(ctx context.Context, req *protos.ListProfilesRequest) (*protos.ListProfilesResponse, error) {
	return nil, nil
}

func (s *server) GetProfile(ctx context.Context, req *protos.GetProfileRequest) (*protos.GetProfileResponse, error) {
	return nil, nil
}

func (s *server) AddProfile(ctx context.Context, req *protos.AddProfileRequest) (*protos.AddProfileResponse, error) {
	if err := s.profilesStore.Add(ctx, req.Profile); err != nil {
		return nil, err
	}

	if err := s.aggStore.Add(ctx, req.Profile); err != nil {
		return nil, err
	}

	return &protos.AddProfileResponse{}, nil
}
