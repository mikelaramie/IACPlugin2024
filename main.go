// Package main converts IaC validation report in JSON to SARIF format.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
    "log"
    "os"
)

const (
	SARIF_SCHEMA                = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json"
	SARIF_VERSION               = "2.1.0"
	IAC_TOOL_DOCUMENTATION_LINK = ""
	IAC_TOOL_NAME               = "analyze-code-security-scc"
)

type violation struct {
	AssetID        string `json:"assetId"`
	NextSteps      string `json:"nextSteps"`
	PolicyID       string `json:"policyId"`
	Severity       string `json:"severity"`
	ViolatedAsset  asset  `json:"violatedAsset"`
	ViolatedPolicy policy `json:"violatedPolicy"`
}

type asset struct {
	Asset     string `json:"asset"`
	AssetType string `json:"assetType"`
}

type policy struct {
	Constraint     string `json:"constraint"`
	ConstraintType string `json:"constraintType"`
}

type sarifOutput struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    tool          `json:"tool"`
	Results []sarifResult `json:"results"`
}

type tool struct {
	Driver driver `json:"driver"`
}

type driver struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	InformationURI string `json:"informationUri"`
	Rules          []rule `json:"rules"`
}

type sarifResult struct {
	RuleID     string         `json:"ruleId"`
	Message    message        `json:"message"`
	Locations  []location     `json:"locations"`
	Properties propertyResult `json:"properties"`
}

type propertyResult struct {
	AssetID   string `json:"assetId"`
	AssetType string `json:"assetType"`
	Asset     string `json:"asset"`
}

type message struct {
	text string `json:"text"`
}

type rule struct {
	ID              string          `json:"id"`
	FullDescription fullDescription `json:"fullDescription"`
	Properties      propertyRule    `json:"properties"`
}

type propertyRule struct {
	Severity            string   `json:"severity"`
	PolicyType          string   `json:"policyType"`
	ComplianceStandard  []string `json:"complianceStandard"`
	PolicySet           string   `json:"policySet"`
	Posture             string   `json:"posture"`
	PostureRevisionID   string   `json:"postureRevisionId"`
	PostureDeploymentID string   `json:"postureDeploymentId"`
	Constraints         string   `json:"constraints"`
	NextSteps           string   `json:"nextSteps"`
}

type fullDescription struct {
	text string `json:"text"`
}

type location struct {
	LogicalLocation []logicalLocation `json:"logicalLocation"`
}

type logicalLocation struct {
	FullyQualifiedName []string `json:"fullyQualifiedName"`
}

type physicalLocation struct {
	ArtifactLocation artifactLocation `json:"artifactLocation"`
}

type artifactLocation struct {
	URI string `json:"uri"`
}

func main() {
	jsonData, err := os.Open("input.json")
	if err != nil {
			log.Fatal(err)
	}
	defer file.Close()

	// Read the file contents
	data, err := ioutil.ReadAll(jsonData)
	if err != nil {
			log.Fatal(err)
	}

	// Decode the JSON data
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
			log.Fatal(err)
	}
	var violations struct {
		Response struct {
			IACValidationReport struct {
				Violations []violation `json:"violations"`
			} `json:"iacValidationReport"`
		} `json:"response"`
	}

	err := json.Unmarshal([]byte(jsonData), &violations)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Convert to SARIF format
	var sarifOutput sarifOutput
	sarifOutput.Schema = SARIF_SCHEMA
	sarifOutput.Version = SARIF_VERSION

	var sarifRuns []sarifRun
	var sarifResults []sarifResult
	var sarifRules []rule

	for _, v := range violations.Response.IACValidationReport.Violations {
		var sarifResult sarifResult
		var sarifRule rule
		var location location
		var logicalLocation logicalLocation
		sarifResult.RuleID = v.PolicyID
		sarifResult.Message.text = fmt.Sprintf("Asset type: %s has a violation, next steps: %s", v.AssetID, v.NextSteps)
		logicalLocation.FullyQualifiedName = []string{v.AssetID}
		location.LogicalLocation = append(location.LogicalLocation, logicalLocation)
		sarifResult.Locations = append(sarifResult.Locations, location)
		sarifResult.Properties = propertyResult{
			AssetID:   v.AssetID,
			AssetType: v.ViolatedAsset.AssetType,
			Asset:     v.ViolatedAsset.Asset,
		}

		sarifRule.ID = v.PolicyID
		sarifRule.FullDescription = fullDescription{text: v.ViolatedPolicy.Constraint}
		sarifRule.Properties = propertyRule{
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

	var sarifRun sarifRun
	sarifRun.Tool.Driver.InformationURI = IAC_TOOL_DOCUMENTATION_LINK
	sarifRun.Tool.Driver.Name = IAC_TOOL_NAME
	sarifRun.Tool.Driver.Version = SARIF_VERSION
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