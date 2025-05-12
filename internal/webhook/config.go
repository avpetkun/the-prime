package webhook

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func LoadConfig(rawConfig string) (cfg Config, err error) {
	if rawConfig == "" {
		return cfg, fmt.Errorf("empty config env CONFIG")
	}
	err = yaml.Unmarshal([]byte(rawConfig), &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, cfg.Validate()
}

type Config struct {
	NatsAddr string `yaml:"nats_addr"`
	HttpPort int    `yaml:"http_port"`
}

func (cfg *Config) Validate() error {
	if cfg.NatsAddr == "" {
		return fmt.Errorf("nats addr is no set")
	}
	if cfg.HttpPort == 0 {
		return fmt.Errorf("http port is no set")
	}
	return nil
}
