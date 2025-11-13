package api

import (
	"context"
	"go-payments-api/internal/application"
	"go-payments-api/internal/infrastructure/api/handler"
	"go-payments-api/internal/settings"
	"go-payments-api/pkg/api"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Application struct {
	BaseApp *application.App
	Server  api.Server[*gin.Engine]

	// Health
	HealthHandler *handler.Health
}

func init() {
	if os.Getenv("ENVIRONMENT") != "test" {
		if err := godotenv.Load(); err != nil {
			panic("Error loading .env file")
		}
	}

	settings.Init()
}

func (a *Application) Start() {
	a.BaseApp.Start(settings.Settings.Metrics.Name)

	a.SetupRoutes()

	ctx := context.Background()
	quitSig := make(chan os.Signal, 1)
	signal.Notify(quitSig, os.Interrupt)

	go func() {
		select {
		case <-quitSig:
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := a.Server.Shutdown(ctx); err != nil {
				a.BaseApp.Logger.Errorf("Server forced to shutdown: %v", err)
			}
			return
		case <-ctx.Done():
			return
		}
	}()

	if err := a.Server.Start(); err != nil {
		a.BaseApp.Logger.Errorf("Failed to start server: %v", err)
	}

	a.BaseApp.Logger.Infof("Server exited properly")
	a.BaseApp.Stop()
}
