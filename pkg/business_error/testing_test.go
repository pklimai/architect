package business_error_test

import (
	"context"
	"errors"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

type fixture struct {
	ctx context.Context
	*assert.Assertions
}

func setUp(t *testing.T) *fixture {
	t.Helper()

	ctrl := minimock.NewController(t)

	return &fixture{
		Assertions: assert.New(ctrl),
	}
}

const (
	errMsgForTest = "error msg for test"
)

var (
	errForTest = errors.New("error for test")
)
