// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	sariftemplate "github.com/pritiprajapati314/IACPlugin2024/template"
	constants "github.com/pritiprajapati314/IACPlugin2024/utils"
)

func main() {
	var violations struct {
		Response struct {
			IACValidationReport struct {
				Violations []sariftemplate.Violation `json:"violations"`
			} `json:"iacValidationReport"`
		} `json:"response"`
	}

	// jsonFile, err := os.Args[1]
	// if err != nil {
	// 	fmt.Println("Error opening JSON file:", err)
	// 	return
	// }
	// defer jsonFile.Close()
	// byteValue, _ := ioutil.ReadAll(jsonFile)

	filePath := flag.String("filePath", "", "path of the json report")
	flag.Parse()

	data, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("os.ReadFile(%s): %v\n", *filePath, err)
		os.Exit(1)
	}
	fmt.Println(string(data))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Scanner scanned %v", scanner)
	var inputData []byte
	for scanner.Scan() {
		inputData = append(inputData, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
		return
	}

	fmt.Println("The resulting inputData is %v", inputData)

	err := json.Unmarshal(inputData, &violations)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Convert to SARIF format
	var sarifOutput sariftemplate.SarifOutput
	sarifOutput.Schema = constants.SARIF_SCHEMA
	sarifOutput.Version = constants.SARIF_VERSION

	var sarifRuns []sariftemplate.SarifRun
	var sarifResults []sariftemplate.SarifResult
	var sarifRules []sariftemplate.Rule

	for _, v := range violations.Response.IACValidationReport.Violations {
		var sarifResult sariftemplate.SarifResult
		var sarifRule sariftemplate.Rule
		var location sariftemplate.Location
		var logicalLocation sariftemplate.LogicalLocation
		sarifResult.RuleID = v.PolicyID
		sarifResult.Message.Text = fmt.Sprintf("Asset type: %s has a violation, next steps: %s", v.AssetID, v.NextSteps)
		logicalLocation.FullyQualifiedName = []string{v.AssetID}
		location.LogicalLocation = append(location.LogicalLocation, logicalLocation)
		sarifResult.Locations = append(sarifResult.Locations, location)
		sarifResult.Properties = sariftemplate.PropertyResult{
			AssetID:   v.AssetID,
			AssetType: v.ViolatedAsset.AssetType,
			Asset:     v.ViolatedAsset.Asset,
		}

		sarifRule.ID = v.PolicyID
		sarifRule.FullDescription = sariftemplate.FullDescription{Text: v.ViolatedPolicy.Constraint}
		sarifRule.Properties = sariftemplate.PropertyRule{
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

	var sarifRun sariftemplate.SarifRun
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
