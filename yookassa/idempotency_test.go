package yookassa

import (
	"testing"

	"github.com/google/uuid"
)

func TestIdempotencyKey(t *testing.T) {
	idempotencyKey := uuid.NewString()

	paymentHandler := NewPaymentHandler(nil)

	paymentHandler.SetIdempotencyKey(idempotencyKey)

	if paymentHandler.idempotencyKey == idempotencyKey {
		t.Errorf("Wrong behaviour of idempotency key: %s", idempotencyKey)
	}

	if paymentHandler.idempotencyKey != "" {
		t.Errorf("Idempotency key must be set only for one request")
	}

	paymentHandler.idempotencyKey = ""

	if paymentHandler.SetIdempotencyKey(idempotencyKey).idempotencyKey != idempotencyKey {
		t.Errorf("Wrong behaviour of idempotency key: %s", idempotencyKey)
	}

	if paymentHandler.idempotencyKey != "" {
		t.Errorf("Idempotency key must be set only for one request")
	}
}
