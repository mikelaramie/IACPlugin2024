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
	"encoding/json"

	"github.com/pritiprajapati314/IACPlugin2024/SARIFConverter/converter"
	"github.com/pritiprajapati314/IACPlugin2024/SARIFConverter/template"
)

var (
	inputFilePath  = flag.String("filePath", "", "path of the input file")
	outputFilePath = flag.String("output", "output.json", "path of the output file")
)

func main() {
	flag.Parse()

	iacReport, err := readAndParseIACScanReport(inputFilePath)
	if err != nil {
		fmt.Printf("readAndParseIACScanReport: %v", err)
		os.Exit(1)
	}

	sarifReport, err := converter.FromIACScanReport(iacReport.Response.IacValidationReport)
	if err != nil {
		fmt.Printf("sarif.FromIACScanReport: %v", err)
		os.Exit(1)
	}

	if err := writeSarifReport(sarifReport, outputFilePath); err != nil {
		fmt.Printf("writeSarifReport(): %v", err)
		os.Exit(1)
	}
}


func readAndParseIACScanReport(filePath *string) (template.IACReportTemplate, error) {
	data, err := os.ReadFile(*filePath)
	if err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	var iacReport template.IACReportTemplate
	if err = json.Unmarshal(data, &iacReport); err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("json.Unmarshal(): %v", err)
	}

	return iacReport, nil
}

func writeSarifReport(sarifReport template.SarifOutput, outputFilePath *string) error {
	sarifJSON, err := json.MarshalIndent(sarifReport, "", "  ")
	if err != nil {
		return fmt.Errorf("json.MarshalIndent: %v", err)
	}

	outputJSON, err := os.Create(*outputFilePath)
	if err != nil {
		return fmt.Errorf("os.Create: %v", err)
	}
	defer outputJSON.Close()

	_, err = outputJSON.Write(sarifJSON)
	if err != nil {
		return fmt.Errorf("outputJSON.Write: %v", err)
	}

	return nil
}