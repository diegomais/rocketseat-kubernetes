package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	if err := runMigrations(db, "create_items_table.sql"); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/data", dataHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	err := db.Ping(context.Background())
	if err != nil {
		http.Error(w, "Database not reachable", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Connection OK"))
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var payload struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil || payload.Name == "" {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		_, err := db.Exec(context.Background(), "INSERT INTO items (name) VALUES ($1)", payload.Name)
		if err != nil {
			http.Error(w, "Failed to insert data", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodGet:
		rows, err := db.Query(context.Background(), "SELECT id, name FROM items")
		if err != nil {
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		type Item struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}
		var items []Item
		for rows.Next() {
			var item Item
			if err := rows.Scan(&item.ID, &item.Name); err == nil {
				items = append(items, item)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func runMigrations(db *pgxpool.Pool, filename string) error {
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}
	_, err = db.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}
	return nil
}
