package settings

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Specification struct {
		Environment string `envconfig:"ENVIRONMENT" default:"dev"`
		HttpServer  HttpServerSpecification
		Metrics     MetricsSpecification
	}

	HttpServerSpecification struct {
		Port         string        `envconfig:"HTTP_SERVER_PORT" default:":8080"`
		ReadTimeout  time.Duration `envconfig:"HTTP_SERVER_READ_TIMEOUT" default:"15s"`
		WriteTimeout time.Duration `envconfig:"HTTP_SERVER_WRITE_TIMEOUT" default:"15s"`
	}

	MetricsSpecification struct {
		Name            string `envconfig:"OTEL_SERVICE_NAME" default:"go-payments-api"`
		Url             string `envconfig:"OTEL_EXPORTER_JAEGER_ENDPOINT" default:"http://localhost:4317"`
		Token           string `envconfig:"SPLUNK_ACCESS_TOKEN"`
		Resource        string `envconfig:"OTEL_RESOURCE_ATTRIBUTES" default:"service.name=go-payments-api"`
		TracesExporter  string `envconfig:"OTEL_TRACES_EXPORTER" default:"jaeger"`
		MetricsExporter string `envconfig:"OTEL_METRICS_EXPORTER" default:""`
	}
)

var Settings Specification

func Init() {
	if err := envconfig.Process("", &Settings); err != nil {
		panic(err.Error())
	}
}

func (s *Specification) IsProduction() bool {
	return s.Environment == "prod"
}

func (s *Specification) IsLocal() bool {
	return s.Environment == "local"
}
