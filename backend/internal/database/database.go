// ด้า
package database

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"kencatexpress/backend/internal/config"
	"kencatexpress/backend/migrations"
)

func Open(cfg config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local&multiStatements=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(2 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Bootstrap(ctx context.Context, db *sql.DB, autoMigrate, autoSeed bool) error {
	entries, err := migrations.Files.ReadDir(".")
	if err != nil {
		return err
	}
	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		if strings.Contains(name, "seed") && !autoSeed {
			continue
		}
		if !strings.Contains(name, "seed") && !autoMigrate {
			continue
		}
		content, err := migrations.Files.ReadFile(name)
		if err != nil {
			return err
		}
		if strings.TrimSpace(string(content)) == "" {
			continue
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("executing migration %s: %w", name, err)
		}
	}
	return nil
}
