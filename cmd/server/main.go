package main

import (
	"context"
	"go-server/internal/app/auth"
	"go-server/internal/app/user"
	"go-server/internal/config"
	"go-server/internal/router"
	"go-server/internal/utils"
	"go-server/pkg/db"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Setup logging
	appConfig := config.LoadAppConfig()
	utils.InitLog(appConfig.Env, appConfig.LogLevel, appConfig.LogFile)
	defer utils.SyncLog()

	// Load DB config and connect
	conn, err := db.Connect(appConfig.DB)
	if err != nil {
		utils.Log.Fatalw("Failed to connect to DB", "error", err)
	}
	defer conn.Close()

	userRepo := user.NewUserRepository(conn)
	authRepo := auth.NewAuthRepository(conn)

	// Setup router
	port := ":" + appConfig.Port
	r := router.NewRouter(conn, appConfig, userRepo, authRepo)

	// Start HTTP server with graceful shutdown
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		utils.Log.Infow("Starting server", "port", appConfig.Port, "env", appConfig.Env, "version", appConfig.Version)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Log.Fatalw("Server error", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	utils.Log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.Log.Fatalw("Server forced to shutdown", "error", err)
	}

	utils.Log.Info("Server exited cleanly.")
}
