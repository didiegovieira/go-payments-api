package entity

import "time"

type PaymentStatus string

const (
	StatusCreated    PaymentStatus = "CREATED"
	StatusProcessing PaymentStatus = "PROCESSING"
	StatusCompleted  PaymentStatus = "COMPLETED"
)

const (
	MethodPix  string = "PIX"
	MethodCard string = "CARD"
)

type Payment struct {
	ID        int64         `json:"id" db:"id"`
	Amount    float64       `json:"amount" db:"amount"`
	Method    string        `json:"method" db:"method"`
	Status    PaymentStatus `json:"status" db:"status"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
}
