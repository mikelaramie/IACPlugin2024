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

	"github.com/pritiprajapati314/IACPlugin2024/ReportValidator/utils"
)

func FetchVoilationFromInputFile(filePath *string) (map[string]int, error) {
	var violationlist utils.IACScanReport

	data, err := os.ReadFile(*filePath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	err = json.Unmarshal(data, &violationlist)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON: %v", err)
	}

	severityCounts := make(map[string]int)

	for _, v := range violationlist.Response.IACValidationReport.Violations {
		severityCounts[strings.ToLower(v.Severity)]++
	}

	return severityCounts, nil
}

func ProcessExpression(expression string) (string, map[string]int, error) {
	pairs := strings.Split(expression, ",")

	if expression == "" {
		return utils.SetDefault()
	}

	var operator = ""
	var userVoilationCount = make(map[string]int)

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		key := strings.ToLower(parts[0])

		if key == "operator" {
			op, err := utils.ValidateOperator(operator, strings.ToLower(parts[1]))
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
			return "", nil, fmt.Errorf("Error converting value to integer: %v", err)
		}

		if value < 0 {
			return "", nil, fmt.Errorf("Validation expression can not have negative values!")
		}

		userVoilationCount[key] = value
	}

	if operator == "" {
		return "", nil, fmt.Errorf("no operator found in expression")
	}

	return operator, userVoilationCount, nil
}
