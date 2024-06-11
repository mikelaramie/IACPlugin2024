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
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IACScanReport struct {
	Response struct {
		IACValidationReport struct {
			Violations []struct {
				Severity string `json:"severity"`
			} `json:"violations"`
		} `json:"iacValidationReport"`
	} `json:"response"`
}

func FetchVoilationFromInputFile(filePath *string) (map[string]int, error) {
	var violationlist IACScanReport

	data, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	err = json.Unmarshal(data, &violationlist)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %v", err)
	}

	severityCounts := make(map[string]int)

	for _, v := range violationlist.Response.IACValidationReport.Violations {
		severityCounts[strings.ToUpper(v.Severity)]++
	}

	return severityCounts, nil
}

func ProcessExpression(expression string) (string, map[string]int, error) {
	pairs := strings.Split(expression, ",")

	if expression == "" {
		return setDefault()
	}

	var operator = ""
	var userVoilationCount = make(map[string]int)

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		key := strings.ToUpper(parts[0])

		if key == "OPERATOR" {
			op, err := validateOperator(operator, strings.ToUpper(parts[1]))
			if err != nil {
				return "", nil, err
			}
			operator = op
			continue
		}

		if _, ok := userVoilationCount[key]; ok {
			return "", nil, fmt.Errorf("duplicate severity found: %v", key)
		}

		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", nil, fmt.Errorf("error converting value to integer: %v", err)
		}

		if value < 0 {
			return "", nil, fmt.Errorf("validation expression can not have negative values")
		}

		userVoilationCount[key] = value
	}

	if operator == "" {
		return "", nil, fmt.Errorf("no operator found in expression")
	}

	return operator, userVoilationCount, nil
}

func validateOperator(finalOperator, expressionOperator string) (string, error) {
	if finalOperator != "" {
		return "", fmt.Errorf("more than one operator found in the expression %v", finalOperator)
	}

	if expressionOperator != "AND" && expressionOperator != "OR" {
		return "", fmt.Errorf("invalid operator: %v", finalOperator)
	}

	return expressionOperator, nil
}

func setDefault() (string, map[string]int, error) {
	userVoilationCount := make(map[string]int)
	userVoilationCount["CRITICAL"] = 1
	userVoilationCount["HIGH"] = 1
	userVoilationCount["MEDIUM"] = 1
	userVoilationCount["LOW"] = 1
	operator := "OR"

	return operator, userVoilationCount, nil
}
