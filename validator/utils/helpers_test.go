package utils

import (
	"testing"
)

func TestValidateOperator(t *testing.T) {
	tests := []struct {
		name               string
		operator           string
		expressionOperator string
		expectedOperator   string
		wantError          bool
	}{
		{
			name:               "Valid single operator (AND)",
			operator:           "",
			expressionOperator: utils.AND,
			expectedOperator:   utils.AND,
			wantError:          false,
		},
		{
			name:               "Valid single operator (OR)",
			operator:           "",
			expressionOperator: utils.OR,
			expectedOperator:   utils.OR,
			wantError:          false,
		},
		{
			name:               "Invalid operator",
			operator:           "",
			expressionOperator: "NOT",
			expectedOperator:   "",
			wantError:          true,
		},
		{
			name:               "operator already set",
			operator:           utils.AND,
			expressionOperator: utils.OR,
			expectedOperator:   "",
			wantError:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			operator, err := utils.ValidateOperator(test.operator, test.expressionOperator)
			if (err != nil) != test.wantError {
				t.Errorf("Expected error: %v, got: %v", test.wantError, err)
			}

			if operator != test.expectedOperator {
				t.Errorf("Expected operator: %v, got: %v", test.expectedOperator, operator)
			}
		})
	}
}

func TestAll(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]bool
		expectedOutput bool
	}{
		{
			name:           "All true",
			input:          map[string]bool{"CRITICAL": true, "HIGH": true, "LOW": true},
			expectedOutput: true,
		},
		{
			name:           "One false",
			input:          map[string]bool{"HIGH": true, "LOW": false, "MEDIUM": true},
			expectedOutput: false,
		},
		{
			name:           "Empty map",
			input:          map[string]bool{},
			expectedOutput: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actualOutput := All(test.input)
			if actualOutput != test.expectedOutput {
				t.Errorf("Expected output: %v, got: %v", test.expectedOutput, actualOutput)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name           string
		input          map[string]bool
		expectedOutput bool
	}{
		{
			name:           "One true",
			input:          map[string]bool{"CRITICAL": true, "HIGH": false, "LOW": false},
			expectedOutput: true,
		},
		{
			name:           "All false",
			input:          map[string]bool{"HIGH": false, "MEDIUM": false, "LOW": false},
			expectedOutput: false,
		},
		{
			name:           "Empty map",
			input:          map[string]bool{},
			expectedOutput: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actualOutput := Any(test.input)
			if actualOutput != test.expectedOutput {
				t.Errorf("Expected output: %v, got: %v", test.expectedOutput, actualOutput)
			}
		})
	}
}
