package business_error_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/zigal0/architect/pkg/business_error"
)

func Test_GetCode(t *testing.T) {
	t.Parallel()

	t.Run("no error", func(t *testing.T) {
		t.Parallel()

		// arrange

		// act
		errorCode := business_error.GetCode(nil)

		// assert
		require.Equal(t, errorCode, business_error.OK)
	})

	t.Run("exist error", func(t *testing.T) {
		t.Parallel()

		// arrange
		err := business_error.New(
			errForTest,
			"no right to operate with this service",
			business_error.PermissionDenied,
		)

		// act
		errorCode := business_error.GetCode(err)

		// assert
		require.Equal(t, errorCode, business_error.PermissionDenied)
	})

	t.Run("exist error", func(t *testing.T) {
		t.Parallel()

		// arrange

		// act
		errorCode := business_error.GetCode(errForTest)

		// assert
		require.Equal(t, errorCode, business_error.Internal)
	})
}
