package gcs

import (
	"fmt"
	"github.com/ducc/profile-collector/protos"
	"github.com/golang/protobuf/ptypes"
	"time"
)

type Formatter func(metadata *protos.ProfileMetadata) (string, error)

func TimeFormatter(metadata *protos.ProfileMetadata) (string, error) {
	startTime, err := ptypes.Timestamp(metadata.StartTime)
	if err != nil {
		return "", err
	}

	endTime, err := ptypes.Timestamp(metadata.EndTime)
	if err != nil {
		return "", err
	}

	startTimeFormatted := startTime.Format(time.RFC3339)
	endTimeFormatted := endTime.Format(time.RFC3339)

	return fmt.Sprintf("%s/%s-%s.profile", metadata.AppName, startTimeFormatted, endTimeFormatted), nil
}
