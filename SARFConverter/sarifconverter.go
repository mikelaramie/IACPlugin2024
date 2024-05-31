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

	fileoperator "github.com/pritiprajapati314/IACPlugin2024/SARFConverter/fileoperator"
	sarif "github.com/pritiprajapati314/IACPlugin2024/SARFConverter/sarif"
)

var (
	inputFilePath  = flag.String("filePath", "", "path of the input file")
	outputFilePath = flag.String("output", "output.json", "path of the output file")
)

func main() {
	flag.Parse()

	iacReport, err := fileoperator.FetchIACScanReport(inputFilePath)
	if err != nil {
		fmt.Printf("fetchIACScanReport: %v", err)
		os.Exit(1)
	}

	sarifReport, err := sarif.GenerateReport(iacReport.Response.IacValidationReport)
	if err != nil {
		fmt.Printf("sarif.GenerateReport: %v", err)
		os.Exit(1)
	}

	fileoperator.ConvertSarifReportToJSONandWriteToOutputFile(sarifReport, outputFilePath)
}