package cache

import "fmt"

type Config struct {
	Addr string `yaml:"addr"`
	Pass string `yaml:"pass"`
	DB   int    `yaml:"db"`
}

func (c *Config) Validate() error {
	if c.Addr == "" {
		return fmt.Errorf("redis addr is no set")
	}
	return nil
}
