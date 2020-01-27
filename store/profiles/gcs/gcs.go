package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/ducc/profile-collector/protos"
	"github.com/golang/protobuf/proto"
)

type store struct {
	bucket    *storage.BucketHandle
	formatter Formatter
}

func New(ctx context.Context, bucket string, formatter Formatter) (*store, error) {
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	b := client.Bucket(bucket)
	if _, err := b.Attrs(ctx); err != nil {
		return nil, err
	}

	return &store{
		bucket:    b,
		formatter: formatter,
	}, nil
}

func (s *store) Add(ctx context.Context, profile *protos.StoredProfile) error {
	data, err := proto.Marshal(profile.Profile)
	if err != nil {
		return err
	}

	objectName, err := s.formatter(profile.Metadata)
	if err != nil {
		return err
	}

	object := s.bucket.Object(objectName)
	writer := object.NewWriter(ctx)

	if _, err := writer.Write(data); err != nil {
		return err
	}

	return nil
}
