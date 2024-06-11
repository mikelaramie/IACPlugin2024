/*
 Copyright 2024 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package evaluate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestComputeViolationState(t *testing.T) {
	tests := []struct {
		name                              string
		severityCounts                    map[string]int
		userViolationCount                map[string]int
		expectedFailureCriteriaViolations map[string]bool
		wantErr                           bool
	}{
		{
			name: "CriticalAndHighLevelSeverityExcceeded",
			severityCounts: map[string]int{
				"CRITICAL": 2,
				"HIGH":     2,
				"MEDIUM":   0,
				"LOW":      0,
			},
			userViolationCount: map[string]int{
				"CRITICAL": 1,
				"HIGH":     1,
				"MEDIUM":   0,
			},
			expectedFailureCriteriaViolations: map[string]bool{
				"CRITICAL": true,
				"HIGH":     true,
				"MEDIUM":   false,
			},
			wantErr: false,
		},
		{
			name: "MediumAndLowLevelSeverityExcceeded",
			severityCounts: map[string]int{
				"CRITICAL": 0,
				"MEDIUM":   3,
				"LOW":      2,
			},
			userViolationCount: map[string]int{
				"CRITICAL": 0,
				"HIGH":     0,
				"MEDIUM":   1,
				"LOW":      1,
			},
			expectedFailureCriteriaViolations: map[string]bool{
				"LOW":      true,
				"MEDIUM":   true,
				"HIGH":     false,
				"CRITICAL": false,
			},
			wantErr: false,
		},
		{
			name: "InvalidSeverity_Error",
			severityCounts: map[string]int{
				"CRITICAL": 2,
				"HIGH":     1,
				"MEDIUM":   0,
			},
			userViolationCount: map[string]int{
				"CRITICAL": 1,
				"HIGH":     2,
				"MEDIUM":   1,
				"invalid":  3,
			},
			expectedFailureCriteriaViolations: nil,
			wantErr:                           true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			failureCriteriaViolations, err := computeViolationState(test.severityCounts, test.userViolationCount)
			if (err != nil) != test.wantErr {
				t.Errorf("Expected error: %v, got: %v", test.wantErr, err)
			}

			if diff := cmp.Diff(test.expectedFailureCriteriaViolations, failureCriteriaViolations); diff != "" {
				t.Errorf("Expected failureCriteriaViolations (+got, -want): %v", diff)
			}
		})
	}
}

func TestIsViolatingSeverity(t *testing.T) {
	tests := []struct {
		name                      string
		operator                  string
		failureCriteriaViolations map[string]bool
		expectedBool              bool
		wantErr                   bool
	}{
		{
			name:     "ANDOperator_SeverityNotViolated",
			operator: "AND",
			failureCriteriaViolations: map[string]bool{
				"MEDIUM": true,
				"HIGH":   false,
			},
			expectedBool: false,
			wantErr:      false,
		},
		{
			name:     "ANDOperator_SeverityViolated",
			operator: "AND",
			failureCriteriaViolations: map[string]bool{
				"MEDIUM": true,
				"HIGH":   true,
			},
			expectedBool: true,
			wantErr:      false,
		},
		{
			name:     "OROperator_SeverityNotViolated",
			operator: "OR",
			failureCriteriaViolations: map[string]bool{
				"key1": false,
				"key2": false,
			},
			expectedBool: false,
			wantErr:      false,
		},
		{
			name:     "OROperator_SeverityViolated",
			operator: "OR",
			failureCriteriaViolations: map[string]bool{
				"key1": true,
				"key2": false,
			},
			expectedBool: true,
			wantErr:      false,
		},
		{
			name:                      "InvalidOperator_Failure",
			operator:                  "RANDOM",
			failureCriteriaViolations: map[string]bool{},
			expectedBool:              true,
			wantErr:                   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			isVoilated, err := isViolatingSeverity(test.operator, test.failureCriteriaViolations)
			if (err != nil) != test.wantErr {
				t.Errorf("Expected error: %v, got: %v", test.wantErr, err)
			}

			if test.expectedBool != isVoilated {
				t.Errorf("Unexpected output want: %v, got: %v", test.expectedBool, isVoilated)
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
			name:           "AllTrue",
			input:          map[string]bool{"CRITICAL": true, "HIGH": true, "LOW": true},
			expectedOutput: true,
		},
		{
			name:           "AllNotTrue",
			input:          map[string]bool{"HIGH": true, "LOW": false, "MEDIUM": true},
			expectedOutput: false,
		},
		{
			name:           "EmptyMap",
			input:          map[string]bool{},
			expectedOutput: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actualOutput := all(test.input)
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
			name:           "OneTrue",
			input:          map[string]bool{"CRITICAL": true, "HIGH": false, "LOW": false},
			expectedOutput: true,
		},
		{
			name:           "AllFalse",
			input:          map[string]bool{"HIGH": false, "MEDIUM": false, "LOW": false},
			expectedOutput: false,
		},
		{
			name:           "EmptyMap",
			input:          map[string]bool{},
			expectedOutput: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actualOutput := any(test.input)
			if actualOutput != test.expectedOutput {
				t.Errorf("Expected output: %v, got: %v", test.expectedOutput, actualOutput)
			}
		})
	}
}
