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

// Package iacTemplate contains the template template for IAC Scan report.
package template

// IACReportTemplate is the template for IAC Scan report.
type IACReportTemplate struct {
	Response Responses `json:"response"`
}

// Response is the response of the IAC Scan report.
type Responses struct {
	Name                string              `json:"name"`
	CreateTime          string              `json:"createTime"`
	UpdateTime          string              `json:"updateTime"`
	IacValidationReport IACValidationReport `json:"iacValidationReport"`
}

type IACValidationReport struct {
	Violations []Violation `json:"violations"`
	Note       string      `json:"note"`
}

// Violation is a violation of the Sarif template.
type Violation struct {
	AssetID         string         `json:"assetId"`
	PolicyID        string         `json:"policyId"`
	ViolatedPosture PostureDetails `json:"violatedPosture"`
	ViolatedPolicy  PolicyDetails  `json:"violatedPolicy"`
	ViolatedAsset   AssetDetails   `json:"violatedAsset"`
	Severity        string         `json:"severity"`
	NextSteps       string         `json:"nextSteps"`
}

type PostureDetails struct {
	PostureDeployment               string `json:"postureDeployment"`
	PostureDeploymentTargetResource string `json:"postureDeploymentTargetResource"`
	Posture                         string `json:"posture"`
	PostureRevisionID               string `json:"postureRevisionId"`
	PolicySet                       string `json:"policySet"`
}

type PolicyDetails struct {
	Constraint          string   `json:"constraint"`
	ConstraintType      string   `json:"constraintType"`
	ComplianceStandards []string `json:"complianceStandards"`
	Description         string   `json:"description"`
}

// AssetDetails is an asset of the Sarif template.
type AssetDetails struct {
	Asset     string `json:"asset"`
	AssetType string `json:"assetType"`
}
