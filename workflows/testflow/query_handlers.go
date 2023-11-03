package testflow

import (
	"go.temporal.io/sdk/workflow"

	"temporalissue/models"
)

const (
	QueryHandlerOrder = "getOrder"
)

func createQueryHandlers(ctx workflow.Context, order *models.Order) error {
	logger := workflow.GetLogger(ctx)

	orderHandler := func() (*models.Order, error) {
		return order, nil
	}

	if err := workflow.SetQueryHandler(ctx, QueryHandlerOrder, orderHandler); err != nil {
		logger.Info("SetQueryHandler failed", "Error", err)
		return err
	}
	return nil
}
