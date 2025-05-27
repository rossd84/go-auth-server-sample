package main

import (
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
	"saas-go-postgres/internal/logger"
	"saas-go-postgres/pkg/db"
)

func main() {
	// Setup logging
	appConfig :=config.LoadAppConfig()
	logger.Init(appConfig.Env, appConfig.LogLevel, appConfig.LogFile)
	defer logger.Sync()

    // Load DB config and connect
    conn, err := db.Connect(appConfig.DB)
    if err != nil {
        logger.Log.Fatalw("Failed to connect to DB", "error", err)
    }
    defer conn.Close()
    user.DB = conn

    // Setup router
    port := ":" + appConfig.Port
	r := router.NewRouter(conn)

    // Start HTTP server with graceful shutdown
    srv := &http.Server{
        Addr:    port,
        Handler: r,
    }

    go func() {
        	logger.Log.Infow("Starting server", "port", appConfig.Port, "env", appConfig.Env, "version", appConfig.Version)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Log.Fatalw("Server error", "error", err)
        }
    }()

    // Graceful shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
    <-quit
    logger.Log.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Log.Fatalw("Server forced to shutdown", "error", err)
    }

    logger.Log.Info("Server exited cleanly.")
}

