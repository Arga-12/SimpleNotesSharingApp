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

func InitDB() *sql.DB {
	psql := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPass, DBName,
	)

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