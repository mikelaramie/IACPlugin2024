package fileoperator

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pritiprajapati314/IACPlugin2024/template"
)

func FetchIACScanReport(filePath *string) (template.IACReportTemplate, error) {
	var iacReport template.IACReportTemplate

	data, err := os.ReadFile(*filePath)
	if err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	if err = json.Unmarshal(data, &iacReport); err != nil {
		return template.IACReportTemplate{}, fmt.Errorf("Error decoding JSON: %v", err)
	}

	return iacReport, nil
}
