package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"

	_ "github.com/lib/pq"
	"saas-go-postgres/internal/router"
	"saas-go-postgres/internal/config"
	"saas-go-postgres/internal/user"
	"saas-go-postgres/pkg/db"
)

func main() {
    // Load DB config and connect
    dbConfig := config.LoadDBConfig()
    conn, err := db.Connect(dbConfig)
    if err != nil {
        log.Fatalf("❌ Failed to connect to DB: %v", err)
    }
    defer conn.Close()
    user.DB = conn

    // Setup router
    port := ":8080"
    r := router.NewRouter(conn)

    // Start HTTP server with graceful shutdown
    srv := &http.Server{
        Addr:    port,
        Handler: r,
    }

    go func() {
        fmt.Printf("🚀 Server running on http://localhost%s\n", port)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server error: %v", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    fmt.Println("\n🛑 Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }

    fmt.Println("✅ Server exited cleanly.")
}

