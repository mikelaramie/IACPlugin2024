package evaluate

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	utils "github.com/pritiprajapati314/IACPlugin2024/validator/utils"
)

func TestProcessExpression(t *testing.T) {
	tests := []struct {
		name                   string
		expression             string
		expectedSeverityCounts map[string]int
		expectedOperator       string
		expectedError          bool
	}{
		{
			name:       "Succeeds",
			expression: "critical:2,high:1,medium:3,operator:or",
			expectedSeverityCounts: map[string]int{
				utils.CRITICAL: 2,
				utils.HIGH:     1,
				utils.MEDIUM:   3,
			},
			expectedOperator: utils.OR,
			expectedError:    false,
		},
		{
			name:                   "duplicate operator present",
			expression:             "critical:2,operator:or,operator:and",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "operator not present",
			expression:             "critical:2,high:1,medium:3",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "duplicate severity present",
			expression:             "critical:2,high:1,medium:3,medium:4,operator:or",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "invalid expression",
			expression:             "critical:invalid,high:1,medium:3",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:       "expression not passed, set default",
			expression: "",
			expectedSeverityCounts: map[string]int{
				utils.CRITICAL: 1,
				utils.HIGH:     1,
				utils.MEDIUM:   1,
				utils.LOW:      1,
			},
			expectedOperator: utils.OR,
			expectedError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			operator, severityCounts, err := processExpression(test.expression)
			if (err != nil) != test.expectedError {
				t.Fatalf("Expected error: %v, got error %v", test.expectedError, err)
			}
			if diff := cmp.Diff(test.expectedSeverityCounts, severityCounts); diff != "" {
				t.Errorf("Expected severityCounts (+got, -want): %v", diff)
			}
			if err == nil && operator != test.expectedOperator {
				t.Errorf("Unexpected operator: expected %v, got %v", test.expectedOperator, operator)
			}
		})
	}
}

func TestComputeVoliationState(t *testing.T) {
	tests := []struct {
		name                              string
		severityCounts                    map[string]int
		userVoilationCount                map[string]int
		expectedFailureCriteriaVoilations map[string]bool
		wantErr                           bool
	}{
		{
			name: "CriticalAndHighLevelSeverityExcceeded",
			severityCounts: map[string]int{
				utils.CRITICAL: 2,
				utils.HIGH:     2,
				utils.MEDIUM:   0,
				utils.LOW:      0,
			},
			userVoilationCount: map[string]int{
				utils.CRITICAL: 1,
				utils.HIGH:     1,
				utils.MEDIUM:   0,
			},
			expectedFailureCriteriaVoilations: map[string]bool{
				utils.CRITICAL: true,
				utils.HIGH:     true,
				utils.MEDIUM:   false,
			},
			wantErr: false,
		},
		{
			name: "MediumAndLowLevelSeverityExcceeded",
			severityCounts: map[string]int{
				utils.CRITICAL: 0,
				utils.MEDIUM:   3,
				utils.LOW:      2,
			},
			userVoilationCount: map[string]int{
				utils.CRITICAL: 0,
				utils.HIGH:     0,
				utils.MEDIUM:   1,
				utils.LOW:      1,
			},
			expectedFailureCriteriaVoilations: map[string]bool{
				utils.LOW:      true,
				utils.MEDIUM:   true,
				utils.HIGH:     false,
				utils.CRITICAL: false,
			},
			wantErr: false,
		},
		{
			name: "InvalidSeverity",
			severityCounts: map[string]int{
				utils.CRITICAL: 2,
				utils.HIGH:     1,
				utils.MEDIUM:   0,
			},
			userVoilationCount: map[string]int{
				utils.CRITICAL: 1,
				utils.HIGH:     2,
				utils.MEDIUM:   1,
				"invalid":      3,
			},
			expectedFailureCriteriaVoilations: nil,
			wantErr:                           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			failureCriteriaVoilations, err := computeVoliationState(test.severityCounts, test.userVoilationCount)
			if (err != nil) != test.wantErr {
				t.Errorf("Expected error: %v, got: %v", test.wantErr, err)
			}

			if diff := cmp.Diff(test.expectedFailureCriteriaVoilations, failureCriteriaVoilations); diff != "" {
				t.Errorf("Expected failureCriteriaVoilations (+got, -want): %v", diff)
			}
		})
	}
}

func TestIsViolatingSeverity(t *testing.T) {
	tests := []struct {
		name                      string
		operator                  string
		failureCriteriaVoilations map[string]bool
		expectedBool              bool
		wantErr                   bool
	}{
		{
			name:     "ANDOperator_SeverityNotViolated",
			operator: utils.AND,
			failureCriteriaVoilations: map[string]bool{
				utils.MEDIUM: true,
				utils.HIGH:   false,
			},
			expectedBool: false,
			wantErr:      false,
		},
		{
			name:     "ANDOperator_SeverityViolated",
			operator: utils.AND,
			failureCriteriaVoilations: map[string]bool{
				utils.MEDIUM: true,
				utils.HIGH:   true,
			},
			expectedBool: true,
			wantErr:      false,
		},
		{
			name:     "OROperator_SeverityNotViolated",
			operator: utils.OR,
			failureCriteriaVoilations: map[string]bool{
				"key1": false,
				"key2": false,
			},
			expectedBool: false,
			wantErr:      false,
		},
		{
			name:     "OROperator_SeverityViolated",
			operator: utils.OR,
			failureCriteriaVoilations: map[string]bool{
				"key1": true,
				"key2": false,
			},
			expectedBool: true,
			wantErr:      false,
		},
		{
			operator:                  "invalid",
			failureCriteriaVoilations: map[string]bool{},
			expectedBool:              true,
			wantErr:                   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			isVoilated, err := isViolatingSeverity(test.operator, test.failureCriteriaVoilations)
			if (err != nil) != test.wantErr {
				t.Errorf("Expected error: %v, got: %v", test.wantErr, err)
			}

			if test.expectedBool != isVoilated {
				t.Errorf("Unexpected output want: %s, got: %v", test.expectedBool, isVoilated)
			}
		})
	}
}
