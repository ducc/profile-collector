package collector

import (
	"context"
	"github.com/ducc/profile-collector/collector/config"
	"github.com/ducc/profile-collector/collector/scraper"
	"github.com/ducc/profile-collector/collector/service"
	"github.com/ducc/profile-collector/protos"
	"github.com/golang/protobuf/ptypes"
	"time"
)

type controller struct {
	conf    config.Config
	scraper scraper.Scraper
	store   protos.StoreClient
}

func Run(ctx context.Context, conf config.Config, scraper scraper.Scraper, store protos.StoreClient) error {
	c := &controller{
		conf:    conf,
		scraper: scraper,
		store:   store,
	}

	c.pollServices(ctx)
	return nil
}

func (c *controller) pollServices(ctx context.Context) {
	for {
		c.scrapeServices(ctx)
		time.Sleep(time.Second * 5) // todo is this the best way to do this
	}
}

func (c *controller) scrapeServices(ctx context.Context) {
	for _, svc := range c.conf.Services() {
		if !svc.ScrapeNow() {
			continue
		}

		c.scrapeService(ctx, svc)
	}
}

func (c *controller) scrapeService(ctx context.Context, svc *service.Service) {
	startTime := time.Now()

	profile, err := c.scraper.Scrape(ctx, svc.Address)
	if err != nil {
		// todo handle error
		return
	}

	endTime := time.Now()
	svc.LastScrape = endTime

	c.storeProfile(ctx, svc, profile, startTime, endTime)
}

func (c *controller) storeProfile(ctx context.Context, svc *service.Service, profile *protos.Profile, startTime, endTime time.Time) {
	timestampStart, err := ptypes.TimestampProto(startTime)
	if err != nil {
		// todo handle error
		return
	}

	timestampEnd, err := ptypes.TimestampProto(endTime)
	if err != nil {
		// todo handle error
		return
	}

	if _, err := c.store.AddProfile(ctx, &protos.AddProfileRequest{
		Profile: &protos.StoredProfile{
			Profile: profile,
			Metadata: &protos.ProfileMetadata{
				AppName:   svc.Name,
				StartTime: timestampStart,
				EndTime:   timestampEnd,
			},
		},
	}); err != nil {
		// todo handle error
	}
}
