package database

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/melkdesousa/gamgo/config"
)

var (
	db     *pgx.Conn
	dbOnce sync.Once
)

func GetDBConnnection() *pgx.Conn {
	dbOnce.Do(func() {
		var err error
		db, err = pgx.Connect(context.Background(), config.MustGetEnv("DB_STRING"))
		if err != nil {
			panic(err)
		}

		// Start a goroutine to periodically check database health
		go startHealthCheck()
	})
	return db
}

// startHealthCheck runs a periodic health check on the database connection
func startHealthCheck() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if db == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := db.Ping(ctx)
		cancel()

		if err != nil {
			// Log the error
			log.Printf("Database health check failed: %v", err)

			// Try to reconnect
			reconnectDB()
		}
	}
}

// reconnectDB attempts to reconnect to the database
func reconnectDB() {
	if db != nil {
		// Close existing connection
		_ = db.Close(context.Background())
	}

	var err error
	db, err = pgx.Connect(context.Background(), config.MustGetEnv("DB_STRING"))
	if err != nil {
		log.Printf("Failed to reconnect to database: %v", err)
	} else {
		log.Println("Successfully reconnected to database")
	}
}
