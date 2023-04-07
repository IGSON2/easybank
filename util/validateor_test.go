package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCurrency struct {
	C string `json:"currency" validate:"required,currency"`
}

func TestValida(t *testing.T) {
	testcase := []struct {
		name         string
		currency     string
		validateFunc func(errs *ErrorResponses) bool
	}{
		{
			name:     "Success",
			currency: "BTC",
			validateFunc: func(errs *ErrorResponses) bool {
				return errs == nil
			},
		},
		{
			name:     "Fail",
			currency: "ADA",
			validateFunc: func(errs *ErrorResponses) bool {
				return errs != nil || len(*errs) > 1
			},
		},
	}

	for _, tc := range testcase {
		TestC := TestCurrency{C: tc.currency}

		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateStruct(TestC)
			passed := tc.validateFunc(errs)
			require.Equal(t, true, passed)

			if errs != nil {
				t.Log(errs.Error())
			}
		})
	}
}
