package testflow_test

import (
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"

	models2 "temporalissue/models"
	"temporalissue/workflows/testflow"
	"temporalissue/workflows/testflow/activities"
	"temporalissue/workflows/testflow/signals"
)

const testOrderID = "E71B3F60-98CA-4FBE-B4BC-E28F226066A0"

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(_, _ string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_OrderCancellation() {
	a := activities.Activities{}
	orderID := testOrderID
	state := &models2.Order{
		ID: strfmt.UUID(orderID),
		OrderCheckout: models2.OrderCheckout{
			Payment: &models2.OrderPayment{
				Code: models2.PaymentCodeCardOnline,
			},
		},
	}

	// mock payment status activity - always new
	s.env.OnActivity(a.PaymentStatus, mock.Anything, mock.Anything).Return(
		&models2.OrderPaymentStatus{Status: models2.PaymentStatusNew}, nil,
	)

	// send order-cancel signal
	s.env.RegisterDelayedCallback(s.cancelSignal(models2.CancelReasonCourierCancelOrder, time.Minute*1))

	s.env.ExecuteWorkflow(testflow.TestFlow, state)
	s.True(s.env.IsWorkflowCompleted())

	// query
	res, err := s.env.QueryWorkflow("getOrder")
	s.NoError(err)
	// get result
	err = res.Get(&state)
	s.NoError(err)
	s.assertOrderStatus(models2.OrderStatusCancel)
}

func (s *UnitTestSuite) getOrder() *models2.Order {
	result := models2.Order{}
	res, err := s.env.QueryWorkflow(testflow.QueryHandlerOrder)
	s.NoError(err)
	// nolint
	err = res.Get(&result)
	return &result
}

func (s *UnitTestSuite) assertOrderStatus(status models2.OrderStatus) {
	result := s.getOrder()
	s.Equal(status, result.Status)
}

func (s *UnitTestSuite) cancelSignal(
	reason models2.CancelReason,
	timer time.Duration,
) (func(), time.Duration) {
	return func() {
		cancellationSignal := signals.SignalCancellation{}
		cancellationSignal.Route = signals.RouteCancellation
		cancellationSignal.Reason = reason
		s.env.SignalWorkflow(signals.RouteCancellation, cancellationSignal)
	}, timer
}
