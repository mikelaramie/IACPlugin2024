// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pritiprajapati314/IACPlugin2024/constants"
	"github.com/pritiprajapati314/IACPlugin2024/sarif"
)

func main() {
	var violations struct {
		Response struct {
			IACValidationReport struct {
				Violations []sarif.Violation `json:"violations"`
			} `json:"iacValidationReport"`
		} `json:"response"`
	}

	jsonFile, err := os.Open("input.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &violations)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Convert to SARIF format
	var sarifOutput sarif.SarifOutput
	sarifOutput.Schema = constants.SARIF_SCHEMA
	sarifOutput.Version = constants.SARIF_VERSION

	var sarifRuns []sarif.SarifRun
	var sarifResults []sarif.SarifResult
	var sarifRules []sarif.Rule

	for _, v := range violations.Response.IACValidationReport.Violations {
		var sarifResult sarif.SarifResult
		var sarifRule sarif.Rule
		var location sarif.Location
		var logicalLocation sarif.LogicalLocation
		sarifResult.RuleID = v.PolicyID
		sarifResult.Message.Text = fmt.Sprintf("Asset type: %s has a violation, next steps: %s", v.AssetID, v.NextSteps)
		logicalLocation.FullyQualifiedName = []string{v.AssetID}
		location.LogicalLocation = append(location.LogicalLocation, logicalLocation)
		sarifResult.Locations = append(sarifResult.Locations, location)
		sarifResult.Properties = sarif.PropertyResult{
			AssetID:   v.AssetID,
			AssetType: v.ViolatedAsset.AssetType,
			Asset:     v.ViolatedAsset.Asset,
		}

		sarifRule.ID = v.PolicyID
		sarifRule.FullDescription = sarif.FullDescription{Text: v.ViolatedPolicy.Constraint}
		sarifRule.Properties = sarif.PropertyRule{
			Severity:   v.Severity,
			PolicyType: v.ViolatedPolicy.ConstraintType,
			// ComplianceStandard:  []string{"STANDARD"},
			// PolicySet:           v.ViolatedPolicy.ConstraintType,
			// Posture:             v.ViolatedPolicy.ConstraintType,
			// PostureRevisionID:   v.ViolatedPolicy.ConstraintType,
			// PostureDeploymentID: v.ViolatedPolicy.ConstraintType,
			Constraints: v.ViolatedPolicy.Constraint,
			NextSteps:   v.NextSteps,
		}

		sarifRules = append(sarifRules, sarifRule)
		sarifResults = append(sarifResults, sarifResult)
	}

	var sarifRun sarif.SarifRun
	sarifRun.Tool.Driver.InformationURI = constants.IAC_TOOL_DOCUMENTATION_LINK
	sarifRun.Tool.Driver.Name = constants.IAC_TOOL_NAME
	sarifRun.Tool.Driver.Version = constants.SARIF_VERSION
	sarifRun.Tool.Driver.Rules = sarifRules
	sarifRun.Results = sarifResults
	sarifRuns = append(sarifRuns, sarifRun)

	sarifOutput.Runs = sarifRuns

	// Convert SARIF output to JSON
	sarifJSON, err := json.MarshalIndent(sarifOutput, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling SARIF:", err)
		return
	}

	// Write SARIF JSON to file
	fmt.Println(string(sarifJSON))

	fmt.Println("Conversion successful. SARIF file generated: output.sarif")
}
