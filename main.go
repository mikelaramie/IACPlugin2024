// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"google3/base/go/flag"

	iacTemplate "google3/experimental/CertifiedOSS/JsonToSarif/iacTemplate"
	sariftemplate "google3/experimental/CertifiedOSS/JsonToSarif/template"
	constants "google3/experimental/CertifiedOSS/JsonToSarif/utils"
)

var (
	inputFilePath  = flag.String("input", "input.json", "path of the input file")
	outputFilePath = flag.String("output", "output.json", "path of the output file")
)

func main() {
	flag.Parse()

	iacReport, err := fetchIACScanReport(inputFilePath)
	if err != nil {
		fmt.Println("Error fetching IAC scan report:", err)
		return
	}

	sarifReport := generateSarifReport(iacReport.Response.IacValidationReport)

	convertSarifReportToJSONandWriteToOutputFile(sarifReport)
}

func fetchIACScanReport(filePath *string) (iacTemplate.IACReportTemplate, error) {
	var iacReport iacTemplate.IACReportTemplate

	data, err := os.ReadFile(*filePath)
	if err != nil {
		return iacTemplate.IACReportTemplate{}, fmt.Errorf("os.ReadFile(%s): %v", *filePath, err)
	}

	if err = json.Unmarshal(data, &iacReport); err != nil {
		return iacTemplate.IACReportTemplate{}, fmt.Errorf("Error decoding JSON: %v", err)
	}

	return iacReport, nil
}

func generateSarifReport(report iacTemplate.IACValidationReport) sariftemplate.SarifOutput {
	policyToViolationMap := getUniqueViolations(report.Violations)
	rules := constructRules(policyToViolationMap)
	results := constructResults(report.Violations)
	sarifReport := sariftemplate.SarifOutput{
		Version: constants.SARIF_VERSION,
		Schema:  constants.SARIF_SCHEMA,
		Runs: []sariftemplate.Run{
			{
				Note: report.Note,
				Tool: sariftemplate.Tool{
					Driver: sariftemplate.Driver{
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
	return sarifReport
}

func getUniqueViolations(violations []iacTemplate.Violation) map[string]iacTemplate.Violation {
	policyToViolationMap := make(map[string]iacTemplate.Violation)
	for _, violation := range violations {
		policyID := violation.PolicyID
		if _, ok := policyToViolationMap[policyID]; !ok {
			policyToViolationMap[policyID] = violation
		}
	}
	return policyToViolationMap
}

func constructRules(policyToViolationMap map[string]iacTemplate.Violation) []sariftemplate.Rule {
	rules := []sariftemplate.Rule{}
	for policyID, violation := range policyToViolationMap {
		rule := sariftemplate.Rule{
			ID: policyID,
			FullDescription: sariftemplate.FullDescription{
				Text: violation.ViolatedPolicy.Description,
			},
			Properties: sariftemplate.RuleProperties{
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
	return rules
}

func constructResults(violations []iacTemplate.Violation) []sariftemplate.Result {
	results := []sariftemplate.Result{}
	for _, violation := range violations {
		result := sariftemplate.Result{
			RuleID: violation.PolicyID,
			Message: sariftemplate.Message{
				Text: fmt.Sprintf("Asset type: %s has a violation, next steps: %s", violation.ViolatedAsset.AssetType, violation.NextSteps),
			},
			Locations: []sariftemplate.Location{
				{
					LogicalLocations: []sariftemplate.LogicalLocation{
						{
							FullyQualifiedName: violation.AssetID,
						},
					},
				},
			},
			Properties: sariftemplate.ResultProperties{
				AssetID:   violation.AssetID,
				Asset:     violation.ViolatedAsset.Asset,
				AssetType: violation.ViolatedAsset.AssetType,
			},
		}
		results = append(results, result)
	}
	return results
}

func convertSarifReportToJSONandWriteToOutputFile(sarifReport sariftemplate.SarifOutput) error {
	sarifJSON, err := json.MarshalIndent(sarifReport, "", "  ")
	if err != nil {
		return fmt.Errorf("Error marshalling SARIF: %v", err)
	}

	outputJSON, err := os.Create(*outputFilePath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer outputJSON.Close()

	_, err = outputJSON.Write(sarifJSON)
	if err != nil {
		return fmt.Errorf("Error writing SARIF JSON to file: %v", err)
	}

	fmt.Println(*outputFilePath)

	return nil
}
