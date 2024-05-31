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

	"github.com/pritiprajapati314/IACPlugin2024/ReportValidator/fileoperator"
	"github.com/pritiprajapati314/IACPlugin2024/ReportValidator/utils"
)

func IsIACScanReportVoilatingSeverity(filePath, expression *string) (bool, error) {
	severityCounts, err := fileoperator.FetchVoilationFromInputFile(filePath)
	if err != nil {
		return false, fmt.Errorf("fetchVoilationFromInputFile failed :%v", err)
	}

	operator, userVoilationCount, err := fileoperator.ProcessExpression(*expression)
	if err != nil {
		return false, fmt.Errorf("processExprssion failed :%v", err)
	}

	failureCriteriaVoilations, err := computeVoliationState(severityCounts, userVoilationCount)
	if err != nil {
		return false, fmt.Errorf("computeVoliationState failed :%v", err)
	}

	return isViolatingSeverity(operator, failureCriteriaVoilations)
}

func computeVoliationState(severityCounts map[string]int, userVoilationCount map[string]int) (map[string]bool, error) {
	failureCriteriaVoilations := make(map[string]bool)

	for k, violationLimit := range userVoilationCount {
		severity := strings.ToLower(k)
		switch severity {
		case utils.CRITICAL:
			if severityCounts[utils.CRITICAL] > violationLimit {
				failureCriteriaVoilations[utils.CRITICAL] = true
			} else {
				failureCriteriaVoilations[utils.CRITICAL] = false
			}
		case utils.HIGH:
			if severityCounts[utils.HIGH] > violationLimit {
				failureCriteriaVoilations[utils.HIGH] = true
			} else {
				failureCriteriaVoilations[utils.HIGH] = false
			}
		case utils.MEDIUM:
			if severityCounts[utils.MEDIUM] > violationLimit {
				failureCriteriaVoilations[utils.MEDIUM] = true
			} else {
				failureCriteriaVoilations[utils.MEDIUM] = false
			}
		case utils.LOW:
			if severityCounts[utils.LOW] > violationLimit {
				failureCriteriaVoilations[utils.LOW] = true
			} else {
				failureCriteriaVoilations[utils.LOW] = false
			}
		default:
			return nil, fmt.Errorf("Invalid severity expression: %v", severity)
		}
	}

	return failureCriteriaVoilations, nil
}

func isViolatingSeverity(operator string, failureCriteriaVoilations map[string]bool) (bool, error) {
	switch operator {
	case utils.AND:
		return utils.All(failureCriteriaVoilations), nil
	case utils.OR:
		return utils.Any(failureCriteriaVoilations), nil
	default:
		return true, fmt.Errorf("Invalid severity operator: %v", operator)
	}
}
