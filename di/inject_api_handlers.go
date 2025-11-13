package di

import (
	"go-payments-api/internal/infrastructure/api/handler"
	"go-payments-api/internal/settings"
	"go-payments-api/pkg/api"
	"go-payments-api/pkg/api/presenter"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var apiHandlersSet = wire.NewSet(
	provideApiServer,
	provideApiPresenter,
	wire.Struct(new(handler.Health), "*"),
)

func provideApiServer() api.Server[*gin.Engine] {
	return api.NewGinServer[*gin.Engine](&http.Server{
		Addr:         settings.Settings.HttpServer.Port,
		ReadTimeout:  settings.Settings.HttpServer.ReadTimeout,
		WriteTimeout: settings.Settings.HttpServer.WriteTimeout,
	})
}

func provideApiPresenter() api.Presenter {
	return presenter.NewJson()
}
