/*
 Copyright 2024 Google LLC

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package sarif

import (
	template "github.com/pritiprajapati314/IACPlugin2024/template"
	constants "github.com/pritiprajapati314/IACPlugin2024/utils"
)

var IACValidationValidReport = template.IACValidationReport{
	Note: "Test Note",
	Violations: []template.Violation{
		template.Violation{
			AssetID:  "Asset 1",
			PolicyID: "P1",
			Severity: "HIGH",
			ViolatedPolicy: template.PolicyDetails{
				Description:         "High-level violation message",
				ConstraintType:      "Type 1",
				ComplianceStandards: []string{"Standard 1"},
			},
			NextSteps: "Next steps 1",
			ViolatedPosture: template.PostureDetails{
				PolicySet:         "Set 1",
				Posture:           "Posture 1",
				PostureRevisionID: "Rev 1",
				PostureDeployment: "Dep 1",
			},
			ViolatedAsset: template.AssetDetails{
				AssetType: "Type 1",
				Asset:     "Asset 1",
			},
		},
	},
}

var IACValidationReportWithInvalidSeverity = template.IACValidationReport{
	Note: "Test Note",
	Violations: []template.Violation{
		template.Violation{
			AssetID:  "Asset 1",
			PolicyID: "P1",
			Severity: "INVALID_SEVERITY",
			ViolatedPolicy: template.PolicyDetails{
				Description:         "High-level violation message",
				ConstraintType:      "Type 1",
				ComplianceStandards: []string{"Standard 1"},
			},
			NextSteps: "Next steps 1",
			ViolatedPosture: template.PostureDetails{
				PolicySet:         "Set 1",
				Posture:           "Posture 1",
				PostureRevisionID: "Rev 1",
				PostureDeployment: "Dep 1",
			},
			ViolatedAsset: template.AssetDetails{
				AssetType: "Type 1",
				Asset:     "Asset 1",
			},
		},
	},
}

var IACValidSarifOutput = template.SarifOutput{
	Version: constants.SARIF_VERSION,
	Schema:  constants.SARIF_SCHEMA,
	Runs: []template.Run{
		{
			Note: "Test Note",
			Tool: template.Tool{
				Driver: template.Driver{
					Name:           constants.IAC_TOOL_NAME,
					Version:        "*******INCORRECT NEEDS TO BE CORRECTED*******",
					InformationURI: constants.IAC_TOOL_DOCUMENTATION_LINK,
					Rules: []template.Rule{
						{
							ID:              "P1",
							FullDescription: template.FullDescription{Text: "High-level violation message"},
							Properties: template.RuleProperties{
								Severity:            "HIGH",
								PolicyType:          "Type 1",
								ComplianceStandard:  []string{"Standard 1"},
								PolicySet:           "Set 1",
								Posture:             "Posture 1",
								PostureRevisionID:   "Rev 1",
								PostureDeploymentID: "Dep 1",
								NextSteps:           "Next steps 1",
							},
						},
					},
				},
			},
			Results: []template.Result{
				{
					RuleID:  "P1",
					Message: template.Message{Text: "Asset type: Type 1 has a violation, next steps: Next steps 1"},
					Locations: []template.Location{
						{
							LogicalLocations: []template.LogicalLocation{
								{FullyQualifiedName: "Asset 1"},
							},
						},
					},
					Properties: template.ResultProperties{
						AssetID:   "Asset 1",
						Asset:     "Asset 1",
						AssetType: "Type 1",
					},
				},
			},
		},
	},
}
