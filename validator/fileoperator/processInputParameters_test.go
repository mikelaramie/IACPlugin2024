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

package fileoperator

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pritiprajapati314/IACPlugin2024/validator/utils"
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
			name:       "ExpressionWithNegativeValue_Failure",
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
			name:                   "DuplicateOperatorPresent_Failure",
			expression:             "critical:2,operator:or,operator:and",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "OperatorNotPresent_Failure",
			expression:             "critical:2,high:1,medium:3",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "DuplicateSeverityPresent_Failure",
			expression:             "critical:2,high:1,medium:3,medium:4,operator:or",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:                   "InvalidExpression_Failure",
			expression:             "critical:invalid,high:1,medium:3",
			expectedSeverityCounts: nil,
			expectedOperator:       "",
			expectedError:          true,
		},
		{
			name:       "ExpressionNotPassed_SetDefault",
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
