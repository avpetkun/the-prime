package worker

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
	Postgres dbx.Config     `yaml:"postgres"`
	Redis    cache.Config   `yaml:"redis"`
	Nats     NatsConfig     `yaml:"nats"`
	Ton      TonConfig      `yaml:"ton"`
	Fragment FragmentConfig `yaml:"fragment"`
	Referral ReferralConfig `yaml:"referral"`
	Workers  WorkersConfig  `yaml:"workers"`
	Bot      tgu.BotConfig  `yaml:"bot"`
	Miniapp  MiniappConfig  `yaml:"miniapp"`
}

type NatsConfig struct {
	Addr string `yaml:"addr"`
}

type TonConfig struct {
	DepositWallet string `yaml:"deposit_wallet"`
}

type FragmentConfig struct {
	WalletSeed string  `yaml:"wallet_seed"`
	AuthCookie string  `yaml:"auth_cookie"`
	MinBalance float64 `yaml:"min_balance"`
	ChatGifts  int64   `yaml:"chatGifts"`
}

type WorkersConfig struct {
	Webhook int `yaml:"webhook"`
	Users   int `yaml:"users"`
	Tasks   int `yaml:"tasks"`
	Chats   int `yaml:"chats"`
	Checks  int `yaml:"checks"`
}

type ReferralConfig struct {
	TotalBonus int64   `yaml:"total_bonus"`
	Levels     []int64 `yaml:"levels"`
}

type MiniappConfig struct {
	WebappURL  string `yaml:"webapp_url"`
	MiniappURL string `yaml:"miniapp_url"`
	SupportURL string `yaml:"support_url"`

	HelloImage string `yaml:"hello_image"`

	En *BotLoc `yaml:"en"`
	Ru *BotLoc `yaml:"ru"`
}

type BotLoc struct {
	WebappButton  string `yaml:"webapp_button"`
	SupportButton string `yaml:"support_button"`
	HelloText     string `yaml:"hello_text"`
}

func (loc *BotLoc) Valid() error {
	if loc.WebappButton == "" {
		return fmt.Errorf("webapp_button is empty")
	}
	if loc.SupportButton == "" {
		return fmt.Errorf("support_button is empty")
	}
	if loc.HelloText == "" {
		return fmt.Errorf("hello_text is empty")
	}
	return nil
}

func (cfg *Config) SetDefaults() {
	cfg.Workers.Webhook = 20
	cfg.Workers.Users = 5
	cfg.Workers.Tasks = 20
	cfg.Workers.Chats = 5
	cfg.Workers.Checks = 5

	cfg.Bot.SetDefaults()

	cfg.Miniapp.WebappURL = "https://miniapp.getprime.me"
	cfg.Miniapp.MiniappURL = "https://t.me/YourPrimeBot/get"
	cfg.Miniapp.SupportURL = "https://t.me/The_Prime_Support_bot"
}

func (cfg *Config) Validate() error {
	if err := cfg.Postgres.Validate(); err != nil {
		return err
	}
	if err := cfg.Redis.Validate(); err != nil {
		return err
	}

	if cfg.Redis.Addr == "" {
		return fmt.Errorf("redis addr is no set")
	}

	if cfg.Nats.Addr == "" {
		return fmt.Errorf("nats addr is no set")
	}

	if cfg.Ton.DepositWallet == "" {
		return fmt.Errorf("ton deposit_wallet address is no set")
	}

	if cfg.Fragment.AuthCookie == "" {
		return fmt.Errorf("fragment auth_cookie is empty")
	}
	if cfg.Fragment.WalletSeed == "" {
		return fmt.Errorf("fragment wallet_seed is empty")
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
	if cfg.Miniapp.HelloImage == "" {
		return fmt.Errorf("bot hello image is empty")
	}
	if err := cfg.Miniapp.En.Valid(); err != nil {
		return fmt.Errorf("loc en: %w", err)
	}
	if err := cfg.Miniapp.Ru.Valid(); err != nil {
		return fmt.Errorf("loc ru: %w", err)
	}

	return nil
}
