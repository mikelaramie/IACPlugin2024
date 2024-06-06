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

// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pritiprajapati314/IACPlugin2024/ReportValidator/evaluate"
)

var (
	filePath   = flag.String("filePath", "", "path of the json file")
	expression = flag.String("expression", "", "condition for validation")
)

func main() {
	flag.Parse()

	isVoilated, err := evaluate.IsIACScanReportVoilatingSeverity(filePath, expression)
	if err != nil {
		fmt.Printf("Failure occured during validation: %v", err)
		os.Exit(1)
	}

	if isVoilated {
		fmt.Printf("Validation Failed!")
		os.Exit(1)
	}

	fmt.Println("Validation Succeeded!")
}