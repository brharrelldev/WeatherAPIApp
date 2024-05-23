package validations

import (
	"errors"
	"github.com/brharrelldev/weatherAPI/constants"
	"github.com/brharrelldev/weatherAPI/models"
	"testing"
)

func TestValidateRequest(t *testing.T) {

	tests := []struct {
		name        string
		req         models.WeatherRequest
		expected    bool
		err         error
		errExpected bool
	}{
		{
			name: "valid city",
			req: models.WeatherRequest{
				City:  "ATLANTA",
				State: "Georgia",
			},
			expected: true,
		},
		{
			name: "invalid city",
			req: models.WeatherRequest{
				City:  "Fake City",
				State: "Minnesota",
			},
			errExpected: true,
			err:         constants.ErrCityInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			resp, err := ValidateRequest(tt.req)
			if (err != nil) == tt.errExpected {
				if !errors.Is(err, tt.err) {
					t.Fatalf("errors expcted %v got %v", tt.errExpected, err)
				}
			}

			if resp != tt.expected {
				t.Fatalf("result expected %v got %v", tt.expected, resp)
			}

		})

	}

}
