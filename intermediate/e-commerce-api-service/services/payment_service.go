package services

import (
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/paymentintent"
)

type PaymentService struct{}

// NewPaymentService creates a new PaymentService with the given Stripe API key.
func NewPaymentService(apiKey string) *PaymentService {
	stripe.Key = apiKey
	return &PaymentService{}
}

// CreatePaymentIntent creates a Stripe PaymentIntent for the given amount.
// The client_secret is returned to the frontend for Stripe.js confirmation.
func (s *PaymentService) CreatePaymentIntent(amountCents int64, currency string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountCents),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	return paymentintent.New(params)
}
