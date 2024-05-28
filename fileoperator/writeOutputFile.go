package fileoperator

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pritiprajapati314/IACPlugin2024/template"
)

func ConvertSarifReportToJSONandWriteToOutputFile(sarifReport template.SarifOutput) error {
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
