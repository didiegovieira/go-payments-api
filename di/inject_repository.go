package di

import (
	"go-payments-api/internal/application/gateway/repository"
	"go-payments-api/internal/infrastructure/database/postgres"
	"log"

	"github.com/google/wire"
)

var repositoriesSet = wire.NewSet(
	ProvidePostgresConnection,
	ProvidePaymentRepository,
)

func ProvidePostgresConnection() (*postgres.DB, func(), error) {
	db, err := postgres.NewConnection()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}

	return db, cleanup, nil
}

func ProvidePaymentRepository(db *postgres.DB) repository.PaymentRepository {
	return postgres.NewPaymentRepository(db.GetConnection())
}
