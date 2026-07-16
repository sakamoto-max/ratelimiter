package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/sakamoto-max/ratelimiter/internal/config"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(ctx context.Context, config *config.Config) error {

	log.Println("starting database migration...")

	url := fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v?sslmode=%v",
		config.Postgres.UserName,
		config.Postgres.Password,
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.Db,
		config.Postgres.SSLmode,
	)

	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres : %w", err)
	}

	defer conn.Close(ctx)

	migrator, err := migrate.NewMigrator(ctx, conn, "SCHEMA_VERSION")
	if err != nil {
		return fmt.Errorf("failed to create migrator : %w", err)
	}

	subTree, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("failed to get the sub tree : %w", err)
	}

	err = migrator.LoadMigrations(subTree)
	if err != nil {
		return fmt.Errorf("failed to load migrations : %w", err)
	}

	currentVersion, err := migrator.GetCurrentVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current version : %w", err)
	}

	err = migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("failed to migrate : %w", err)
	}

	if currentVersion == int32(len(migrator.Migrations)) {
		log.Printf("database is up to date : version %v", currentVersion)
	} else {
		log.Printf("database has migrated from %v to %v", currentVersion, len(migrator.Migrations))
	}

	return nil
}
