package testcnts

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestEnvironment() {
	os.Setenv("DOCKER_HOST", "unix:///Users/akartikov/.lima/avito/sock/docker.sock")
	os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")
}

func StartPostgresContainer(ctx context.Context) (testcontainers.Container, string, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", err
	}

	// Получаем параметры подключения
	host, _ := postgresContainer.Host(ctx)
	port, _ := postgresContainer.MappedPort(ctx, "5432/tcp")

	dsn := fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())

	// Ожидаем готовности базы данных перед применением миграций
	if err := waitForDatabase(dsn); err != nil {
		return nil, "", err
	}

	if err := applyMigrations(dsn); err != nil {
		return nil, "", err
	}

	return postgresContainer, dsn, nil
}

func waitForDatabase(dsn string) error {
	for i := 0; i < 30; i++ {
		db, err := sqlx.Connect("postgres", dsn)
		if err == nil {
			if err := db.Ping(); err == nil {
				db.Close()
				return nil
			}
			db.Close()
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("database not ready after 30 seconds")
}

func applyMigrations(dsn string) error {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	// Найдем корень проекта по файлу go.mod
	migrationsPath, err := findMigrationsPath()
	if err != nil {
		return err
	}

	goose.SetTableName("goose_db_version")
	return goose.Up(db.DB, migrationsPath)
}

func findMigrationsPath() (string, error) {
	rootDir, err := findProjectRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(rootDir, "db", "migrations"), nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Ищем go.mod файл, поднимаясь по директориям
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find project root with go.mod")
}
