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

package utils

import "fmt"

func ValidateOperator(finalOperator, expressionOperator string) (string, error) {
	if finalOperator != "" {
		return "", fmt.Errorf("more than one operator found in the expression %v", finalOperator)
	}

	switch expressionOperator {
	case AND:
		return AND, nil
	case OR:
		return OR, nil
	default:
		return "", fmt.Errorf("Invalid operator: %v", finalOperator)
	}
}

func All(failureCriteriaVoilations map[string]bool) bool {
	if len(failureCriteriaVoilations) == 0 {
		return false
	}

	for _, v := range failureCriteriaVoilations {
		if !v {
			return false
		}
	}
	return true
}

func Any(failureCriteriaVoilations map[string]bool) bool {
	if len(failureCriteriaVoilations) == 0 {
		return false
	}

	for _, v := range failureCriteriaVoilations {
		if v {
			return true
		}
	}
	return false
}

func SetDefault() (string, map[string]int, error) {
	userVoilationCount := make(map[string]int)
	userVoilationCount[CRITICAL] = DefaultCritical
	userVoilationCount[HIGH] = DefaultHigh
	userVoilationCount[MEDIUM] = DefaultMedium
	userVoilationCount[LOW] = DefaultLow
	operator := DefaultOperator

	return operator, userVoilationCount, nil
}
