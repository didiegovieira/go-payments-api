package handler

import (
	"go-payments-api/internal/application/dto"
	"go-payments-api/internal/application/usecase"
	"go-payments-api/pkg/api"
	appErr "go-payments-api/pkg/errors"
	"go-payments-api/pkg/metrics"
	"go-payments-api/pkg/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type CreatePayment struct {
	UseCase   usecase.CreatePayment
	Presenter api.Presenter
}

// CreatePayment godoc
// @Summary      Create a new payment
// @Description  Create a new payment and publish event to Kafka
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        payment body dto.CreatePaymentInput true "Payment data"
// @Success      201  {object}  dto.CreatePaymentOutput
// @Failure      400  {object}  api.HttpError
// @Failure      500  {object}  api.HttpError
// @Router       /payments [post]
func (h *CreatePayment) Handle() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		reqCtx, span := metrics.StartSpan(ctx.Request.Context(), "CreatePaymentHandler.Handle")
		defer span.End()

		var input dto.CreatePaymentInput
		if err := ctx.ShouldBindJSON(&input); err != nil {
			metrics.AddSpanEvent(reqCtx, "bind.failed", attribute.String("error", err.Error()))
			h.Presenter.Error(ctx, appErr.HttpBadRequest("Invalid request body"))
			return
		}

		// Validate input
		if err := validator.ValidateStruct(input); err != nil {
			metrics.AddSpanEvent(reqCtx, "validation.failed", attribute.String("error", err.Error()))
			h.Presenter.Error(ctx, appErr.HttpBadRequest(err.Error()))
			return
		}

		// Execute use case
		output, err := h.UseCase.Execute(reqCtx, input)
		if err != nil {
			metrics.AddSpanEvent(reqCtx, "usecase.failed", attribute.String("error", err.Error()))
			h.Presenter.Error(ctx, err)
			return
		}

		metrics.AddSpanAttributes(reqCtx, attribute.Int64("payment.created.id", output.ID))
		h.Presenter.Present(ctx, output, http.StatusCreated)
	}
}
