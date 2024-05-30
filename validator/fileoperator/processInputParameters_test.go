package fileoperator

import (
	"reflect"
	"testing"

	"github.com/pritiprajapati314/IACPlugin2024/validator/utils"
)

func TestProcessExpression(t *testing.T) {
	tests := []struct {
		name                   string
		expression             string
		ExpectedOperator       string
		expectedVoilationCount map[string]int
		wantError              bool
	}{
		{
			name:             "Valid expression with operator",
			expression:       "AND:CRITICAL:2,HIGH:3",
			expectedOperator: utils.AND,
			expectedVoilationCount: map[string]int{
				utils.CRITICAL: 2,
				utils.HIGH:     3,
			},
			wantError: false,
		},
		{
			name:             "Valid expression with OR operator",
			expression:       "OR:HIGH:0,LOW:1",
			ExpectedOperator: utils.OR,
			expectedVoilationCount: map[string]int{
				utils.HIGH: 0,
				utils.LOW:  1,
			},
			wantError: false,
		},
		{
			name:                   "Empty expression",
			expression:             "",
			ExpectedOperator:       utils.SetDefault(),
			expectedVoilationCount: map[string]int{},
			wantError:              false,
		},
		{
			name:                   "Invalid operator",
			expression:             "INVALID:CRITICAL:2",
			ExpectedOperator:       "",
			expectedVoilationCount: nil,
			wantError:              true,
		},
		{
			name:                   "Duplicate severity",
			expression:             "AND:CRITICAL:2,CRITICAL:3",
			ExpectedOperator:       "",
			expectedVoilationCount: nil,
			wantError:              true,
		},
		{
			name:                   "Invalid violation limit",
			expression:             "AND:CRITICAL:abc,HIGH:3",
			ExpectedOperator:       "",
			expectedVoilationCount: nil,
			wantError:              true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			operator, violationCount, err := ProcessExpression(test.expression)

			if (err != nil) != test.wantError {
				t.Errorf("Expected error: %v, got: %v", test.wantError, err)
			}

			if operator != test.ExpectedOperator {
				t.Errorf("Expected operator: %v, got: %v", test.ExpectedOperator, operator)
			}

			if !reflect.DeepEqual(violationCount, test.expectedVoilationCount) {
				t.Errorf("Expected violation count: %v, got: %v", test.expectedVoilationCount, violationCount)
			}
		})
	}
}
