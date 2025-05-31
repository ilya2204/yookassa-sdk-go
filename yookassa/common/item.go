package yoocommon

import (
	"encoding/json"
	"unicode/utf8"
)

type Item struct {
	// parameter with the name of the product or service
	Description string `json:"description"`

	// parameter with the amount per unit of product
	Quantity int `json:"quantity"`

	// parameter specifying the quantity of goods (only integers, for example 1)
	Amount *Amount `json:"amount"`

	// parameter with the fixed value 1 (price without VAT)
	VatCode int `json:"vat_code"`

	Measure string `json:"measure,omitempty"`

	PaymentSubject string `json:"payment_subject,omitempty"`

	PaymentMode string `json:"payment_mode,omitempty"`
}

const MAX_DESCRIPTION_LENGTH = 128

func (u Item) MarshalJSON() ([]byte, error) {
	type Alias Item

	truncated := u

	if utf8.RuneCountInString(u.Description) > MAX_DESCRIPTION_LENGTH {
		runes := []rune(u.Description)
		truncated.Description = string(runes[:MAX_DESCRIPTION_LENGTH])
	}

	return json.Marshal(Alias(truncated))
}
