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

	utils "github.com/pritiprajapati314/IACPlugin2024/ReportValidator/utils"
)

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
			name: "InvalidSeverity_Error",
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
			name:                      "InvalidOperator_Failure",
			operator:                  "RANDOM",
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
				t.Errorf("Unexpected output want: %v, got: %v", test.expectedBool, isVoilated)
			}
		})
	}
}
