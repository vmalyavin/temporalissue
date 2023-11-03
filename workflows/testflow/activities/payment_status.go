package activities

import (
	"context"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	models2 "temporalissue/models"
)

type Activities struct {
}

var PaymentStatusOptions = workflow.ActivityOptions{
	StartToCloseTimeout: time.Second * 5,
	RetryPolicy: &temporal.RetryPolicy{
		MaximumAttempts: 1,
	},
}

func WithPaymentStatusOptions(ctx workflow.Context) workflow.Context {
	return workflow.WithActivityOptions(ctx, PaymentStatusOptions)
}

// PaymentStatus - stub activity to get payment status from external service
func (a *Activities) PaymentStatus(_ context.Context, _ *models2.Order) (*models2.OrderPaymentStatus, error) {
	return nil, nil
}

func PaymentStatusExecute(ctx workflow.Context, order *models2.Order) (*models2.OrderPaymentStatus, error) {
	var paymentStatus models2.OrderPaymentStatus
	if err := workflow.ExecuteActivity(
		WithPaymentStatusOptions(ctx),
		(&Activities{}).PaymentStatus,
		order,
	).Get(ctx, &paymentStatus); err != nil {
		return nil, err
	}
	return &paymentStatus, nil
}
