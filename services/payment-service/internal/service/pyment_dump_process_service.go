package service

import (
	"context"
	"log"
	"ride-sharing/services/payment-service/internal/domain"
	"ride-sharing/services/payment-service/pkg/types"

	"github.com/stripe/stripe-go/v81"
)

type dumpClient struct {
	config *types.PaymentConfig
}

func NewDumpClient(config *types.PaymentConfig) domain.PaymentProcessorServiceService {
	stripe.Key = config.StripeSecretKey

	return &dumpClient{
		config: config,
	}
}

func (s *dumpClient) CreatePaymentSession(ctx context.Context, amount int64, currency string, metadata map[string]string) (string, error) {
	log.Print(
		amount, currency, metadata,
	)
	return "id", nil
}
