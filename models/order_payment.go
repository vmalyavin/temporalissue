package models

type OrderPayment struct {
	ID     string        `json:"id,omitempty"`
	Code   PaymentCode   `json:"code"`
	Status PaymentStatus `json:"status,omitempty"`
}

type OrderPaymentStatus struct {
	ID     string        `json:"ID"`
	Status PaymentStatus `json:"status"`
}

type PaymentCode int64

const (
	PaymentCodeCardOnline PaymentCode = 3
)

type PaymentStatus int64

const (
	PaymentStatusNew PaymentStatus = 0
)
