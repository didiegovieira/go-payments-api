package repository

import (
	"context"
	"go-payments-api/internal/domain/entity"
)

type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id int64) (*entity.Payment, error)
}
