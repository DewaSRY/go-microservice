package domain

import (
	"context"

	"ride-sharing/services/payment-service/pkg/types"
)

type PaymentHandlerService interface {
	CreatePaymentSession(ctx context.Context, tripID, userID, driverID string, amount int64, currency string) (*types.PaymentIntent, error)
}

type PaymentProcessorServiceService interface {
	CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error)
}
