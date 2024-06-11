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
	"fmt"
	"strings"

	"github.com/mikelaramie/IACPlugin2024/ReportValidator/fileoperator"
)

func IsIACScanReportViolatingSeverity(filePath, expression *string) (bool, error) {
	severityCounts, err := fileoperator.FetchViolationFromInputFile(filePath)
	if err != nil {
		return false, fmt.Errorf("fetchViolationFromInputFile hi failed :%v", err)
	}

	operator, userViolationCount, err := fileoperator.ProcessExpression(*expression)
	if err != nil {
		return false, fmt.Errorf("processExpression failed :%v", err)
	}

	failureCriteriaViolations, err := computeViolationState(severityCounts, userViolationCount)
	if err != nil {
		return false, fmt.Errorf("computeViolationState failed :%v", err)
	}

	return isViolatingSeverity(operator, failureCriteriaViolations)
}

func computeViolationState(severityCounts map[string]int, userViolationCount map[string]int) (map[string]bool, error) {
	failureCriteriaViolations := make(map[string]bool)

	for k, violationLimit := range userViolationCount {
		severity := strings.ToUpper(k)
		switch severity {
		case "CRITICAL":
			if severityCounts["CRITICAL"] >= violationLimit {
				failureCriteriaViolations["CRITICAL"] = true
			} else {
				failureCriteriaViolations["CRITICAL"] = false
			}
		case "HIGH":
			if severityCounts["HIGH"] >= violationLimit {
				failureCriteriaViolations["HIGH"] = true
			} else {
				failureCriteriaViolations["HIGH"] = false
			}
		case "MEDIUM":
			if severityCounts["MEDIUM"] >= violationLimit {
				failureCriteriaViolations["MEDIUM"] = true
			} else {
				failureCriteriaViolations["MEDIUM"] = false
			}
		case "LOW":
			if severityCounts["LOW"] >= violationLimit {
				failureCriteriaViolations["LOW"] = true
			} else {
				failureCriteriaViolations["LOW"] = false
			}
		default:
			return nil, fmt.Errorf("invalid severity expression: %v", severity)
		}
	}

	return failureCriteriaViolations, nil
}

func isViolatingSeverity(operator string, failureCriteriaViolations map[string]bool) (bool, error) {
	switch operator {
	case "AND":
		return all(failureCriteriaViolations), nil
	case "OR":
		return any(failureCriteriaViolations), nil
	default:
		return true, fmt.Errorf("invalid severity operator: %v", operator)
	}
}

func all(failureCriteriaViolations map[string]bool) bool {
	if len(failureCriteriaViolations) == 0 {
		return false
	}

	for _, v := range failureCriteriaViolations {
		if !v {
			return false
		}
	}
	return true
}

func any(failureCriteriaViolations map[string]bool) bool {
	if len(failureCriteriaViolations) == 0 {
		return false
	}

	for _, v := range failureCriteriaViolations {
		if v {
			return true
		}
	}
	return false
}
