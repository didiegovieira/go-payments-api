package test

import (
	"go-payments-api/internal/settings"
	"go-payments-api/pkg/http"
	"go-payments-api/pkg/log"
	"go-payments-api/pkg/metrics"
	"io"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"testing"

	log2 "go-payments-api/pkg/log/implement"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/mock/gomock"
)

func Setup(t *testing.T, setEnv *map[string]string) *gomock.Controller {
	log.Logger = log2.Discard()
	metrics.Tracer = trace.NewNoopTracerProvider().Tracer("test")

	t.Setenv("ENVIRONMENT", "local")

	if setEnv != nil {
		for k, v := range *setEnv {
			t.Setenv(k, v)
		}
	}

	settings.Init()

	return gomock.NewController(t)
}

func Data(file string) io.Reader {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "..")
	f, err := os.Open(path.Join(dir, "test/data/"+file+".json"))
	if err != nil {
		panic(err)
	}
	return f
}

func FatalIfErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func MockHttpWrapperResponse(status int, testDataFile string, headers map[string]string) http.Response {
	body, _ := ioutil.ReadAll(Data(testDataFile))

	return http.Response{
		StatusCode: status,
		Body:       body,
		Headers:    headers,
	}
}

func MockHttpWrapperInvalidResponse(status int, headers map[string]string) http.Response {
	return http.Response{
		StatusCode: status,
		Body:       []byte("invalid json"),
		Headers:    headers,
	}
}
