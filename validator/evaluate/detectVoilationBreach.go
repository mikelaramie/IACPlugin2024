package evaluate

import (
	"fmt"
	"strings"

	"github.com/pritiprajapati314/IACPlugin2024/jsonprocessor/fileoperator"
	"github.com/pritiprajapati314/IACPlugin2024/jsonprocessor/utils"
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
