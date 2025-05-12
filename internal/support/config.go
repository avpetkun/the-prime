package support

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
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
}

func (cfg *Config) Validate() error {
	if cfg.BotToken == "" {
		return fmt.Errorf("bot token is no set")
	}
	if cfg.ChatID == 0 {
		return fmt.Errorf("chat id is no set")
	}
	return nil
}
