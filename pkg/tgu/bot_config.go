package tgu

import "fmt"

type BotConfig struct {
	Token string `yaml:"token"`

	RateMsgLimit int `yaml:"rate_msg_limit"`
	RateApiLimit int `yaml:"rate_api_limit"`
}

func (cfg *BotConfig) SetDefaults() {
	cfg.RateMsgLimit = 25
	cfg.RateApiLimit = 25
}

func (cfg *BotConfig) Validate() error {
	if cfg.Token == "" {
		return fmt.Errorf("bot token is empty")
	}
	if cfg.RateMsgLimit <= 0 {
		return fmt.Errorf("invalid rate bot msg limit")
	}
	if cfg.RateApiLimit <= 0 {
		return fmt.Errorf("invalid rate bot api limit")
	}
	return nil
}
