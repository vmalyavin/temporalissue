package testflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	models "temporalissue/models"
	"temporalissue/workflows/testflow/activities"
)

func TestFlow(ctx workflow.Context, order *models.Order) error {
	if err := createQueryHandlers(ctx, order); err != nil {
		return workflow.ErrCanceled
	}

	// after wf cancel - poll for payment status 10m
	defer func() {
		if order.Status != models.OrderStatusCancel || order.Payment.Status != models.PaymentStatusNew {
			return
		}
		newCtx, _ := workflow.NewDisconnectedContext(ctx)
		order.Payment.Status = PollPaymentStatus10Min(newCtx, order)
	}()

	// wait f
	order.Payment.Status = PollPaymentStatus20MinOrReceiveCancel(ctx, order)
	if order.Status == models.OrderStatusCancel {
		return workflow.ErrCanceled
	}

	return nil
}

// pollPaymentStatus makes timer for pereodic polling every pollinterval
func pollPaymentStatus(
	ctx workflow.Context,
	pollInterval time.Duration,
	dst *models.PaymentStatus,
	order *models.Order,
) (workflow.Future, func(workflow.Future)) {
	return workflow.NewTimer(ctx, pollInterval),
		func(_ workflow.Future) {
			logger := workflow.GetLogger(ctx)
			paymentStatus, err := activities.PaymentStatusExecute(ctx, order)
			if err != nil {
				logger.Error("Failed get payment status", "Error", err)
				return
			}
			*dst = paymentStatus.Status
		}
}
