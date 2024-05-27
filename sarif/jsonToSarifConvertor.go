package sarif

import (
	"fmt"

	"github.com/pritiprajapati314/IACPlugin2024/template"
)

func GenerateReport(report template.IACValidationReport) (template.SarifOutput, error) {
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

func validateSeverity(severity string) bool {
	if severity != constants.CRITICAL && severity != constants.HIGH && severity != constants.MEDIUM && severity != constants.LOW {
		return false
	}
	return true
}
