package testflow

import (
	"time"

	"go.temporal.io/sdk/workflow"

	"temporalissue/models"
)

func PollPaymentStatus10Min(
	ctx workflow.Context, order *models.Order,
) models.PaymentStatus {
	var status models.PaymentStatus

	// cycle over timer - 10min
	itsOver := false
	overTimer := workflow.NewTimer(ctx, 10*time.Minute)
	overCallback := func(_ workflow.Future) { itsOver = true }

	for !itsOver {
		pollTimer, pollCallback := pollPaymentStatus(
			ctx, 2*time.Minute, &status, order,
		)
		selector := workflow.NewSelector(ctx).
			AddFuture(pollTimer, pollCallback).
			AddFuture(overTimer, overCallback)
		// no cancel receive
		selector.Select(ctx)

		// end cycle if paymentstatus is changed from new
		if status != models.PaymentStatusNew {
			break
		}
	}

	return status
}
