// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"flag"
	"fmt"
	"os"

	fileoperator "github.com/pritiprajapati314/IACPlugin2024/fileoperator"
	sarif "github.com/pritiprajapati314/IACPlugin2024/sarif"
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

	return fileoperator.ConvertSarifReportToJSONandWriteToOutputFile(sarifReport, outputFilePath)
}
