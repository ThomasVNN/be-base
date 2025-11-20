package transport_test

import (
	"context"
	"errors"
	"testing"

	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/log"
	"git-codecommit.ap-southeast-1.amazonaws.com/v1/repos/be-base/transport"
)

func TestLogErrorHandler(t *testing.T) {
	var output []interface{}

	logger := log.Logger(log.LoggerFunc(func(keyvals ...interface{}) error {
		output = append(output, keyvals...)
		return nil
	}))

	errorHandler := transport.NewLogErrorHandler(logger)

	err := errors.New("error")

	errorHandler.Handle(context.Background(), err)

	if output[1] != err {
		t.Errorf("expected an error log event: have %v, want %v", output[1], err)
	}
}
