package api

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/avpetkun/the-prime/internal/cache"
	"github.com/avpetkun/the-prime/internal/dbx"
	"github.com/avpetkun/the-prime/pkg/tgu"
)

func LoadConfig(rawConfig string) (cfg Config, err error) {
	if rawConfig == "" {
		return cfg, fmt.Errorf("empty config env CONFIG")
	}
	cfg.SetDefaults()
	err = yaml.Unmarshal([]byte(rawConfig), &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, cfg.Validate()
}

type Config struct {
	Postgres dbx.Config    `yaml:"postgres"`
	Redis    cache.Config  `yaml:"redis"`
	Nats     NatsConfig    `yaml:"nats"`
	HTTP     HTTPConfig    `yaml:"http"`
	Bot      tgu.BotConfig `yaml:"bot"`
	Miniapp  MiniappConfig `yaml:"miniapp"`
	Ton      TonConfig     `yaml:"ton"`
}

type NatsConfig struct {
	Addr string `yaml:"addr"`
}

type HTTPConfig struct {
	PublicPort  int `yaml:"public_port"`
	PrivatePort int `yaml:"private_port"`
}

type MiniappConfig struct {
	WebappURL  string `yaml:"webapp_url"`
	MiniappURL string `yaml:"miniapp_url"`
	SupportURL string `yaml:"support_url"`

	InviteImage string `yaml:"invite_image"`

	En *BotLoc `yaml:"en"`
	Ru *BotLoc `yaml:"ru"`
}

type BotLoc struct {
	InviteText   string `yaml:"invite_text"`
	InviteButton string `yaml:"invite_button"`
}

func (loc *BotLoc) Valid() error {
	if loc.InviteText == "" {
		return fmt.Errorf("invite_text is empty")
	}
	if loc.InviteButton == "" {
		return fmt.Errorf("invite_button is empty")
	}
	return nil
}

type TonConfig struct {
	DepositWallet string `yaml:"deposit_wallet"`
	CommentPrefix string `yaml:"comment_prefix"`
}

func (cfg *Config) SetDefaults() {
	cfg.Bot.SetDefaults()

	cfg.Miniapp.WebappURL = "https://miniapp.getprime.me"
	cfg.Miniapp.MiniappURL = "https://t.me/YourPrimeBot/get"
	cfg.Miniapp.SupportURL = "https://t.me/The_Prime_Support_bot"

	cfg.HTTP.PublicPort = 8091
	cfg.HTTP.PrivatePort = 8092
}

func (cfg *Config) Validate() error {
	if err := cfg.Postgres.Validate(); err != nil {
		return err
	}
	if err := cfg.Redis.Validate(); err != nil {
		return err
	}

	if err := cfg.Bot.Validate(); err != nil {
		return err
	}

	if cfg.Miniapp.WebappURL == "" {
		return fmt.Errorf("bot webapp url is empty")
	}
	if cfg.Miniapp.MiniappURL == "" {
		return fmt.Errorf("bot miniapp url is empty")
	}
	if cfg.Miniapp.SupportURL == "" {
		return fmt.Errorf("bot support url is empty")
	}
	if cfg.Miniapp.InviteImage == "" {
		return fmt.Errorf("bot invite image is empty")
	}
	if err := cfg.Miniapp.En.Valid(); err != nil {
		return fmt.Errorf("loc en: %w", err)
	}
	if err := cfg.Miniapp.Ru.Valid(); err != nil {
		return fmt.Errorf("loc ru: %w", err)
	}

	if cfg.Ton.DepositWallet == "" {
		return fmt.Errorf("ton deposit wallet is empty")
	}

	if cfg.HTTP.PublicPort == 0 || cfg.HTTP.PrivatePort == 0 {
		return fmt.Errorf("invalid http ports")
	}

	return nil
}
