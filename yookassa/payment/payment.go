// Package yoopayment describes all the necessary entities for working with YooMoney Payments.
package yoopayment

import (
	"encoding/json"
	"fmt"
	"time"

	yoocommon "github.com/ilya2204/yookassa-sdk-go/yookassa/common"
)

// The Payment object contains all currently relevant information
// about the payment. The object is generated during creation of a payment,
// and sent in response to any payment-related requests.
type Payment struct {
	// Payment ID in YooMoney.
	ID string `json:"id,omitempty"`

	// Payment Status. Possible values: pending, waiting_for_capture, succeeded, and canceled.
	Status Status `json:"status,omitempty"`

	// Payment Amount. Sometimes YooMoney's partners charge additional
	// commission from the users that is not included in this amount.
	Amount *yoocommon.Amount `json:"amount,omitempty"`

	// Amount of payment to be received by the store: the amount value minus the YooMoney commission.
	IncomeAmount *yoocommon.Amount `json:"income_amount,omitempty"`

	// Capture defines automatic acceptance of payment
	Capture bool `json:"capture,omitempty"`

	// Description of the transaction (maximum 128 characters) displayed in your YooMoney
	// Merchant Profile, and shown to the user during checkout. For example,
	// "Payment for order No. 72 for user@yoomoney.ru".
	Description string `json:"description,omitempty" binding:"max=128"`

	// Payment Recipient.
	Recipient *Recipient `json:"recipient,omitempty"`

	// Payment Receipt
	Receipt *Receipt `json:"receipt,omitempty"`

	// Payment method used for this payment.
	PaymentMethod PaymentMethoder `json:"payment_method,omitempty"`

	// ID of the saved payment method
	PaymentMethodID string `json:"payment_method_id,omitempty"`

	// SavePaymentMethod indicates whether the payment method should be saved
	SavePaymentMethod bool `json:"save_payment_method,omitempty"`

	// Time of order creation, based on UTC and specified in the ISO 8601 format.
	// Example: 2017-11-03T11:52:31.827Z
	CapturedAt *time.Time `json:"captured_at,omitempty"`

	// Time of order creation, based on UTC and specified in the ISO 8601 format.
	// Example: 2017-11-03T11:52:31.827Z
	CreatedAt *time.Time `json:"created_at,omitempty"`

	// The period during which you can cancel or capture a payment for free.
	// The payment with the waiting_for_capture status will be automatically
	// canceled at the specified time. Based on UTC and specified in the ISO 8601 format.
	// Example: 2017-11-03T11:52:31.827Z
	ExpiresAt *time.Time `json:"expires_at,omitempty"`

	// Selected payment confirmation scenario.
	// For payments requiring confirmation from the user.
	Confirmation Confirmer `json:"confirmation,omitempty"`

	// The attribute of a test transaction.
	Test bool `json:"test,omitempty"`

	// The amount refunded to the user. Specified if the payment has successful refunds.
	RefundedAmount *yoocommon.Amount `json:"refunded_amount,omitempty"`

	// The attribute of a paid order.
	Paid bool `json:"paid,omitempty"`

	// Availability of the option to make a refund via API.
	Refundable bool `json:"refundable,omitempty"`

	// Status of receipt delivery.
	ReceiptRegistration Status `json:"receipt_registration,omitempty"`

	// Any additional data you might require for processing payments
	// (for example, your internal order ID), specified as a “key-value” pair and
	// returned in response from YooMoney. Limitations: no more than 16 keys,
	// no more than 32 characters in the key’s title, no more than 512 characters
	// in the key’s value, data type is a string in the UTF-8 format.
	Metadata interface{} `json:"metadata,omitempty"`

	// Commentary to the canceled status: who and why canceled the payment.
	CancellationDetails *yoocommon.CancellationDetails `json:"cancellation_details,omitempty"`

	// Payment authorization details.
	AuthorizationDetails *AuthorizationDetails `json:"authorization_details,omitempty"`

	// Information about money distribution: the amounts of transfers and
	// the stores to be transferred to.
	Transfers []Transfer `json:"transfers,omitempty"`

	// The deal within which the payment is being carried out.
	Deal *Deal `json:"deal,omitempty"`

	// The identifier of the customer in your system, such as email address or phone number.
	// No more than 200 characters.
	MerchantCustomerID string `json:"merchant_customer_id,omitempty" binding:"max=200"`
}

func (p *Payment) GetConfirmationToken() (string, error) {
	m, ok := p.Confirmation.(map[string]any)
	if !ok {
		return "", fmt.Errorf("confirmation is not a map")
	}
	raw, ok := m["confirmation_token"]
	if !ok {
		return "", fmt.Errorf("confirmation_token not found in confirmation map")
	}
	token, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("confirmation_token is not a string, got %T", raw)
	}
	if token == "" {
		return "", fmt.Errorf("confirmation_token is empty")
	}
	return token, nil
}

func (p *Payment) GetInvoiceIdFromMetadata() (string, error) {
	m, ok := p.Metadata.(map[string]any)
	if !ok {
		return "", fmt.Errorf("metadata is not a map[string]any, got %T", p.Metadata)
	}
	raw, ok := m["invoice_id"]
	if !ok {
		return "", fmt.Errorf("invoice_id not found in metadata map")
	}
	id, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("invoice_id is not a string, got %T", raw)
	}
	if id == "" {
		return "", fmt.Errorf("invoice_id is empty")
	}
	return id, nil
}

func (p *Payment) GetBasePaymentMethod() (BasePaymentMethod, error) {
	return convertPaymentMethod[BasePaymentMethod](p.PaymentMethod)
}

func (p *Payment) GetPaymentMethodWithCard() (PaymentMethodWithCard, error) {
	return convertPaymentMethod[PaymentMethodWithCard](p.PaymentMethod)
}

func (p *Payment) GetPaymentMethodSbp() (SBP, error) {
	sbp, err := convertPaymentMethod[SBP](p.PaymentMethod)

	if err != nil {
		return sbp, err
	}

	fmt.Println("ALE", sbp.Type)

	if sbp.paymentMethod.Type != PaymentTypeSBP {
		return sbp, fmt.Errorf("payment method is not SBP, got %s", sbp.paymentMethod.Type)
	}

	return sbp, nil
}

func convertPaymentMethod[T any](pm interface{}) (T, error) {
	var zero T

	switch v := pm.(type) {
	case *T:
		return *v, nil
	case T:
		return v, nil
	case map[string]interface{}:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return zero, fmt.Errorf("failed to marshal map: %w", err)
		}

		var result T
		if err := json.Unmarshal(jsonData, &result); err != nil {
			return zero, fmt.Errorf("failed to unmarshal to %T: %w", result, err)
		}

		return result, nil
	default:
		return zero, fmt.Errorf("unsupported payment method type: %T", pm)
	}
}
