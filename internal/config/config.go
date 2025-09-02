package config

import (
	"fmt"

	"github.com/joho/godotenv"
)

type PgConfigSearcher interface {
	Get() (*PgConfig, error)
}

func Load(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	return nil
}

type PgConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DbName   string
}

func (cfg *PgConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}
