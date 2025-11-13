package dto

import "time"

type CreatePaymentInput struct {
	Amount float64 `json:"amount" binding:"required,gt=0" example:"100.50"`
	Method string  `json:"method" binding:"required,oneof=PIX CARD" example:"PIX"`
}

type CreatePaymentOutput struct {
	ID        int64     `json:"id" example:"1"`
	Amount    float64   `json:"amount" example:"100.50"`
	Method    string    `json:"method" example:"PIX"`
	Status    string    `json:"status" example:"CREATED"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T10:00:00Z"`
}

type PaymentEvent struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	EventType string    `json:"event_type"`
}
