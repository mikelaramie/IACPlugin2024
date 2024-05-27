// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	template "github.com/pritiprajapati314/IACPlugin2024/template"
	constants "github.com/pritiprajapati314/IACPlugin2024/utils"
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

	sarifReport, err := generateSarifReport(iacReport.Response.IacValidationReport)
	if err != nil {
		fmt.Printf("FenerateSarifReport: %v", err)
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

func generateSarifReport(report template.IACValidationReport) (template.SarifOutput, error) {
	policyToViolationMap := getUniqueViolations(report.Violations)
	rules, err := constructRules(policyToViolationMap)
	if err != nil {
		return template.SarifOutput{}, fmt.Errorf("constructRules: %v", err)
	}

	results := constructResults(report.Violations)
	sarifReport := template.SarifOutput{
		Version: constants.SARIF_VERSION,
		Schema:  constants.SARIF_SCHEMA,
		Runs: []template.Run{
			{
				Note: report.Note,
				Tool: template.Tool{
					Driver: template.Driver{
						Name:           constants.IAC_TOOL_NAME,
						Version:        "*******INCORRECT NEEDS TO BE CORRECTED*******",
						InformationURI: constants.IAC_TOOL_DOCUMENTATION_LINK,
						Rules:          rules,
					},
				},
				Results: results,
			},
		},
	}
	return sarifReport, nil
}

func getUniqueViolations(violations []template.Violation) map[string]template.Violation {
	policyToViolationMap := make(map[string]template.Violation)
	for _, violation := range violations {
		policyID := violation.PolicyID
		if _, ok := policyToViolationMap[policyID]; !ok {
			policyToViolationMap[policyID] = violation
		}
	}

	return policyToViolationMap
}

func constructRules(policyToViolationMap map[string]template.Violation) ([]template.Rule, error) {
	rules := []template.Rule{}
	for policyID, violation := range policyToViolationMap {
		if !validateSeverity(violation.Severity) {
			return nil, fmt.Errorf("validateSeverity: %s", violation.Severity)
		}

		rule := template.Rule{
			ID: policyID,
			FullDescription: template.FullDescription{
				Text: violation.ViolatedPolicy.Description,
			},
			Properties: template.RuleProperties{
				Severity:            violation.Severity,
				PolicyType:          violation.ViolatedPolicy.ConstraintType,
				ComplianceStandard:  violation.ViolatedPolicy.ComplianceStandards,
				PolicySet:           violation.ViolatedPosture.PolicySet,
				Posture:             violation.ViolatedPosture.Posture,
				PostureRevisionID:   violation.ViolatedPosture.PostureRevisionID,
				PostureDeploymentID: violation.ViolatedPosture.PostureDeployment,
				Constraints:         violation.ViolatedPolicy.Constraint,
				NextSteps:           violation.NextSteps,
			},
		}

		rules = append(rules, rule)
	}
	return rules, nil
}

func constructResults(violations []template.Violation) []template.Result {
	results := []template.Result{}
	for _, violation := range violations {
		result := template.Result{
			RuleID: violation.PolicyID,
			Message: template.Message{
				Text: fmt.Sprintf("Asset type: %s has a violation, next steps: %s", violation.ViolatedAsset.AssetType, violation.NextSteps),
			},
			Locations: []template.Location{
				{
					LogicalLocations: []template.LogicalLocation{
						{
							FullyQualifiedName: violation.AssetID,
						},
					},
				},
			},
			Properties: template.ResultProperties{
				AssetID:   violation.AssetID,
				Asset:     violation.ViolatedAsset.Asset,
				AssetType: violation.ViolatedAsset.AssetType,
			},
		}
		results = append(results, result)
	}
	return results
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

func validateSeverity(severity string) bool {
	if severity != constants.CRITICAL && severity != constants.HIGH && severity != constants.MEDIUM && severity != constants.LOW {
		return false
	}
	return true
}
