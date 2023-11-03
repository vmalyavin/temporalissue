package testflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"temporalissue/models"
	"temporalissue/workflows/testflow/signals"
)

func PollPaymentStatus20MinOrReceiveCancel(
	ctx workflow.Context, order *models.Order,
) models.PaymentStatus {
	var status models.PaymentStatus

	// cycle over timer - 20min
	itsOver := false
	overTimer := workflow.NewTimer(ctx, 20*time.Minute)
	overCallback := func(_ workflow.Future) { itsOver = true }

	for !itsOver {
		pollTimer, pollCallback := pollPaymentStatus(
			ctx, 3*time.Minute, &status, order,
		)
		selector := workflow.NewSelector(ctx).
			AddFuture(pollTimer, pollCallback).
			AddFuture(overTimer, overCallback).
			AddReceive(signals.CancellationReceive(ctx, order))
		selector.Select(ctx)

		// end cycle if paymentstatus is changed from new; or order is cancelled
		if status == models.PaymentStatusNew || order.Status == models.OrderStatusCancel {
			break
		}
	}

	return status
}
