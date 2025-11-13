package di

import (
	"go-payments-api/internal/application/usecase"

	"github.com/google/wire"
)

var provideCreatePaymentUseCase = wire.NewSet(
	usecase.NewCreatePaymentUseCase,
	wire.Bind(new(usecase.CreatePayment), new(*usecase.CreatePaymentImplementation)),
)

var usecasesSet = wire.NewSet(
	provideCreatePaymentUseCase,
)
