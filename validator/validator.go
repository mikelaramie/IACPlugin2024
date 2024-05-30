// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pritiprajapati314/IACPlugin2024/validator/evaluate"
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
