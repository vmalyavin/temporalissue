package models

import (
	"github.com/go-openapi/strfmt"
)

type Order struct {
	OrderCheckout

	ID     strfmt.UUID `json:"id"`
	Status OrderStatus `json:"status"`
}

type OrderCheckout struct {
	Payment *OrderPayment `json:"payment"`
}

type OrderStatus int32

const (
	OrderStatusNew    OrderStatus = 0
	OrderStatusCancel OrderStatus = 99
)

type CancelReason int32

const (
	CancelReasonCourierCancelOrder CancelReason = 8
)
