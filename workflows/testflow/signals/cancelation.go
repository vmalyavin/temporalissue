package signals

import (
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/workflow"

	"temporalissue/models"
)

const RouteCancellation = "RouteCancellation"

type SignalCancellation struct {
	Route string

	Reason models.CancelReason
	Author string
}

func CancellationReceive(ctx workflow.Context, order *models.Order) (workflow.ReceiveChannel, func(workflow.ReceiveChannel, bool)) {
	return workflow.GetSignalChannel(ctx, RouteCancellation), onCancellationReceive(ctx, order)
}

func onCancellationReceive(ctx workflow.Context, order *models.Order) func(c workflow.ReceiveChannel, _ bool) {
	return func(c workflow.ReceiveChannel, _ bool) {
		logger := workflow.GetLogger(ctx)
		signal, err := ReceiveSignal[SignalCancellation](ctx, c)
		if err != nil {
			logger.Error("failed parse signal", "Error", err)
			return
		}
		CancelOrder(ctx, order, signal.Reason)
	}
}

func CancelOrder(_ workflow.Context, order *models.Order, _ models.CancelReason) {
	order.Status = models.OrderStatusCancel
}

func ReceiveSignal[T any](ctx workflow.Context, c workflow.ReceiveChannel) (*T, error) {

	var signal interface{}
	c.Receive(ctx, &signal)

	var message T
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			ToTimeHookFunc()),
		Result: &message,
	})

	if err != nil {
		return nil, err
	}

	if err := decoder.Decode(signal); err != nil {
		return nil, err
	}
	return &message, nil
}

func ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
		// Convert it by parsing
	}
}
