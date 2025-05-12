package dbx

import "fmt"

type Config struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Addr string `yaml:"addr"`
	Name string `yaml:"name"`
}

func (cfg *Config) Validate() error {
	switch {
	case cfg.Addr == "":
		return fmt.Errorf("postgres addr is no set")
	case cfg.User == "":
		return fmt.Errorf("postgres user is no set")
	case cfg.Name == "":
		return fmt.Errorf("postgres name is no set")
	default:
		return nil
	}
}
