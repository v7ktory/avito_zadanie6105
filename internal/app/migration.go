//go:build migrate

package app

import (
	"errors"
	"log"
	"log/slog"
	sl "log/slog"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"

	// migrate tools
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

func init() {
	databaseURL, ok := os.LookupEnv("POSTGRES_CONN")
	if !ok || len(databaseURL) == 0 {
		log.Fatalf("migrate: environment variable not declared: POSTGRES_CONN")
	}

	databaseURL += "?sslmode=disable"

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", databaseURL)
		if err == nil {
			break
		}

		sl.Debug("Migrate: pgdb is trying to connect", slog.Any("attempts left", attempts))
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: pgdb connect error: %s", err)
	}

	err = m.Up()
	defer func() { _, _ = m.Close() }()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		sl.Debug("Migrate: no change")
		return
	}

	sl.Debug("Migrate: up success")
}
