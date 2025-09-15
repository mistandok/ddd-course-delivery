package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type PgConfigSearcher interface {
	Get() (*PgConfig, error)
}

type HttpConfigSearcher interface {
	Get() (*HttpConfig, error)
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
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
}

type HttpConfig struct {
	Host string
	Port int
}

func (cfg *HttpConfig) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

type envHttpConfigSearcher struct{}

func NewHttpConfigSearcher() HttpConfigSearcher {
	return &envHttpConfigSearcher{}
}

func (e *envHttpConfigSearcher) Get() (*HttpConfig, error) {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		host = "localhost"
	}

	portStr := os.Getenv("HTTP_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}

	return &HttpConfig{
		Host: host,
		Port: port,
	}, nil
}
