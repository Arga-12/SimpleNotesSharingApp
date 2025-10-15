package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/lib/pq" // masih pakai pq
)

// DBConfig holds database connection configuration
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// LoadDBConfig loads database configuration from environment variables
func LoadDBConfig() *DBConfig {
	return &DBConfig{
		Host:     getenv("DB_HOST", "localhost"),
		Port:     getenv("DB_PORT", "5432"),
		User:     getenv("DB_USER", "postgres"),
		Password: getenv("DB_PASSWORD", "postgres"),
		Name:     getenv("DB_NAME", "authdb"),
	}
}

// GetConnectionString returns PostgreSQL connection string from config
func (c *DBConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func getenv(k, fallback string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return fallback
}

func InitDB() *sql.DB {
	config := LoadDBConfig()
	psql := config.GetConnectionString()

	db, err := sql.Open("postgres", psql)
	if err != nil {
		log.Fatalf("open db : %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("ping db: %v", err)
	}

	// Jalankan migrasi SQL otomatis
	if err := runMigrations(db); err != nil {
		log.Fatalf("gagal migrasi: %v", err)
	}

	return db
}

func runMigrations(db *sql.DB) error {
	ctx := context.Background()

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_version (
			version TEXT PRIMARY KEY
		);
	`)
	if err != nil {
		return err
	}

	files, err := os.ReadDir("migrations")
	if err != nil {
		return err
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	applied := map[string]bool{}
	rows, _ := db.QueryContext(ctx, "SELECT version FROM schema_version")
	for rows.Next() {
		var v string
		rows.Scan(&v)
		applied[v] = true
	}
	rows.Close()

	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".sql") || applied[f.Name()] {
			continue
		}

		content, err := os.ReadFile("migrations/" + f.Name())
		if err != nil {
			return err
		}

		log.Printf("Menjalankan migrasi cihuyyy %s...", f.Name())
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("error di %s: %w", f.Name(), err)
		}

		_, err = db.ExecContext(ctx, "INSERT INTO schema_version (version) VALUES ($1)", f.Name())
		if err != nil {
			return err
		}
	}

	log.Println("Migrasi selesai, cihuuyyy!")
	return nil
}
