package di

import (
	"go-payments-api/internal/application"
	log "go-payments-api/pkg/log/implement"

	"github.com/google/wire"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var commonSet = wire.NewSet(
	provideLogger,
	provideTracer,
	gatewaysSet,
	wire.Struct(new(application.App), "*"),
)

func provideLogger() log.Logger {
	return log.NewLogrus()
}

func provideTracer() trace.Tracer {
	return otel.Tracer("")
}
