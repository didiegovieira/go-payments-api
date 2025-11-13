package usecase

import (
	"context"
	"fmt"
	"go-payments-api/internal/application/dto"
	"go-payments-api/internal/application/gateway/repository"
	"go-payments-api/internal/domain/entity"
	"go-payments-api/internal/infrastructure/messaging/kafka"
	"go-payments-api/pkg/base"
	"go-payments-api/pkg/metrics"
	"log"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
)

type CreatePayment = base.UseCase[dto.CreatePaymentInput, *dto.CreatePaymentOutput]

type CreatePaymentImplementation struct {
	repository repository.PaymentRepository
	publisher  kafka.Publisher
}

func NewCreatePaymentUseCase(
	repository repository.PaymentRepository,
	publisher kafka.Publisher,
) *CreatePaymentImplementation {
	return &CreatePaymentImplementation{
		repository: repository,
		publisher:  publisher,
	}
}

func (uc *CreatePaymentImplementation) Execute(ctx context.Context, input dto.CreatePaymentInput) (*dto.CreatePaymentOutput, error) {
	ctx, span := metrics.StartSpan(ctx, "CreatePaymentUseCase.Execute")
	defer span.End()

	log.Printf("🔵 Starting payment creation - Amount: %.2f, Method: %s", input.Amount, input.Method)

	metrics.AddSpanAttributes(ctx,
		attribute.Float64("payment.amount", input.Amount),
		attribute.String("payment.method", input.Method),
	)

	// Validate payment method
	if input.Method != entity.MethodPix && input.Method != entity.MethodCard {
		log.Printf("❌ Invalid payment method: %s", input.Method)
		return nil, fmt.Errorf("invalid payment method: %s", input.Method)
	}

	// Create payment entity
	payment := &entity.Payment{
		Amount: input.Amount,
		Method: input.Method,
	}

	// Save to database
	log.Printf("💾 Saving payment to database...")
	if err := uc.repository.Create(ctx, payment); err != nil {
		log.Printf("❌ Failed to save payment to database: %v", err)
		metrics.AddSpanEvent(ctx, "payment.creation.failed", attribute.String("error", err.Error()))
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	log.Printf("✅ Payment saved to database with ID: %d", payment.ID)
	metrics.AddSpanAttributes(ctx, attribute.Int64("payment.id", payment.ID))

	// Publish event to Kafka
	event := dto.PaymentEvent{
		ID:        payment.ID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    string(payment.Status),
		CreatedAt: payment.CreatedAt,
		EventType: "payment.created",
	}

	log.Printf("📤 Publishing event to Kafka - Topic: payment.events, Key: %d", payment.ID)
	if err := uc.publisher.Publish(ctx, "payment.events", strconv.FormatInt(payment.ID, 10), event); err != nil {
		log.Printf("❌ Failed to publish event to Kafka: %v", err)
		metrics.AddSpanEvent(ctx, "kafka.publish.failed", attribute.String("error", err.Error()))
		// Don't fail the request, just log
	} else {
		log.Printf("✅ Event published successfully to Kafka")
	}

	// Return output
	return &dto.CreatePaymentOutput{
		ID:        payment.ID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    string(payment.Status),
		CreatedAt: payment.CreatedAt,
	}, nil
}
