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
	for _, v := range failureCriteriaVoilations {
		if !v {
			return false
		}
	}
	return true
}

func Any(failureCriteriaVoilations map[string]bool) bool {
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
