// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	sarif "github.com/pritiprajapati314/IACPlugin2024/sarif"
	template "github.com/pritiprajapati314/IACPlugin2024/template"
)

var (
	inputFilePath  = flag.String("filePath", "", "path of the input file")
	outputFilePath = flag.String("output", "output.json", "path of the output file")
)

func main() {
	iacReport, err := fetchIACScanReport(inputFilePath)
	if err != nil {
		fmt.Printf("FetchIACScanReport: %v", err)
		os.Exit(1)
	}

	sarifReport, err := sarif.GenerateReport(iacReport.Response.IacValidationReport)
	if err != nil {
		fmt.Printf("sarif.GenerateReport: %v", err)
		os.Exit(1)
	}

	convertSarifReportToJSONandWriteToOutputFile(sarifReport)
}

func fetchIACScanReport(filePath *string) (template.IACReportTemplate, error) {
	var iacReport template.IACReportTemplate

	data, err := os.ReadFile(*filePath)
	if err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	if err = json.Unmarshal(data, &iacReport); err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("json.Unmarshal: %v", err)
	}

	return iacReport, nil
}

func convertSarifReportToJSONandWriteToOutputFile(sarifReport template.SarifOutput) error {
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

	fmt.Println(*outputFilePath)

	return nil
}
