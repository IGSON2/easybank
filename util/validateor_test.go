package util

import (
	"errors"
	"fmt"
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
		validateFunc func(errs []*ErrorResponse) bool
	}{
		{
			name:     "Success",
			currency: "BTC",
			validateFunc: func(errs []*ErrorResponse) bool {
				return errs == nil
			},
		},
		{
			name:     "Fail",
			currency: "ADA",
			validateFunc: func(errs []*ErrorResponse) bool {
				return errs != nil || len(errs) > 1
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
				err := combineErrors(errs)
				t.Log(err)
			}
		})
	}
}

func combineErrors(errs []*ErrorResponse) error {
	var errorString string
	for i, er := range errs {
		errorString += fmt.Sprintf("%02d - FailedField:%v, Tag:%v, Value:%v", i+1, er.FailedField, er.Tag, er.Value)
	}
	return errors.New(errorString)
}
