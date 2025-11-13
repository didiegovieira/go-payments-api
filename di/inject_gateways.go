package di

import (
	"go-payments-api/pkg/http"

	"github.com/google/wire"
)

var gatewaysSet = wire.NewSet(
	http.NewWrapper,
	wire.Bind(new(http.Wrapper), new(*http.WrapperImpl)),
)
