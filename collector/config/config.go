package config

import "github.com/ducc/profile-collector/collector/service"

type Config interface {
	Load() error
	Services() []*service.Service
}

type config struct {
	path string
}

func New(path string) (Config, error) {
	c := &config{
		path: path,
	}

	return c, c.Load()
}

func (c *config) Load() error {
	// todo
	return nil
}

func (c *config) Services() []*service.Service {
	// todo
	return nil
}
