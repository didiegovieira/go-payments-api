package postgres

import (
	"context"
	"database/sql"
	"go-payments-api/internal/application/gateway/repository"
	"go-payments-api/internal/domain/entity"
	"time"

	_ "github.com/lib/pq"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) repository.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	query := `
        INSERT INTO payments (amount, method, status, created_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `

	payment.CreatedAt = time.Now()
	payment.Status = entity.StatusCreated

	err := r.db.QueryRowContext(
		ctx,
		query,
		payment.Amount,
		payment.Method,
		payment.Status,
		payment.CreatedAt,
	).Scan(&payment.ID)

	return err
}

func (r *paymentRepository) FindByID(ctx context.Context, id int64) (*entity.Payment, error) {
	query := `
        SELECT id, amount, method, status, created_at
        FROM payments
        WHERE id = $1
    `

	payment := &entity.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.Amount,
		&payment.Method,
		&payment.Status,
		&payment.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return payment, err
}
