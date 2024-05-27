package main

import (
	"fmt"
	"reflect"
	"testing"

	template "google3/experimental/CertifiedOSS/JsonToSarif/template"
)

func TestGetUniqueViolations(t *testing.T) {
	testCases := []struct {
		name     string
		input    []template.Violation
		expected map[string]template.Violation
	}{
		{
			name:     "No violations",
			input:    []template.Violation{},
			expected: map[string]template.Violation{},
		},
		{
			name: "Multiple unique violations",
			input: []template.Violation{
				{PolicyID: "policy1", Severity: "violation1"},
				{PolicyID: "policy2", Severity: "violation2"},
				{PolicyID: "policy3", Severity: "violation3"},
			},
			expected: map[string]template.Violation{
				"policy1": {PolicyID: "policy1", Severity: "violation1"},
				"policy2": {PolicyID: "policy2", Severity: "violation2"},
				"policy3": {PolicyID: "policy3", Severity: "violation3"},
			},
		},
		{
			name: "Duplicate violations",
			input: []template.Violation{
				{PolicyID: "policy1", Severity: "violation1"},
				{PolicyID: "policy2", Severity: "violation2"},
				{PolicyID: "policy1", Severity: "voilation1"},
			},
			expected: map[string]template.Violation{
				"policy1": {PolicyID: "policy1", Severity: "violation1"},
				"policy2": {PolicyID: "policy2", Severity: "violation2"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getUniqueViolations(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				fmt.Printf("Expected %v, got %v\n", tc.expected, result)
			}
		})
	}
}

func TestConstructRules(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]template.Violation
		expected []template.Rule
	}{
		{
			name:     "Empty Input",
			input:    map[string]template.Violation{},
			expected: []template.Rule{},
		},
		{
			name: "Multiple Violations",
			input: map[string]template.Violation{
				"policy1": {
					PolicyID:        "policy1",
					Severity:        "HIGH",
					ViolatedPolicy:  template.PolicyDetails{Description: "Description 1", ConstraintType: "Type 1", ComplianceStandards: []string{"Standard 1"}},
					ViolatedPosture: template.PostureDetails{PolicySet: "Set 1", Posture: "Posture 1", PostureRevisionID: "Rev 1", PostureDeployment: "Dep 1"},
					NextSteps:       "Next steps 1",
				},
				"policy2": {
					PolicyID:       "policy2",
					Severity:       "MEDIUM",
					ViolatedPolicy: template.PolicyDetails{Description: "Description 2", ConstraintType: "Type 2"},
					NextSteps:      "Next steps 2",
				},
			},
			expected: []template.Rule{
				{
					ID:              "policy1",
					FullDescription: template.FullDescription{Text: "Description 1"},
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
				{
					ID:              "policy2",
					FullDescription: template.FullDescription{Text: "Description 2"},
					Properties: template.RuleProperties{
						Severity:   "MEDIUM",
						PolicyType: "Type 2",
						NextSteps:  "Next steps 2",
					},
				},
			},
		},
		{
			name: "Missing Fields",
			input: map[string]template.Violation{
				"policy3": {
					PolicyID:        "policy3",
					Severity:        "LOW",
					ViolatedPolicy:  template.PolicyDetails{},
					ViolatedPosture: template.PostureDetails{},
					NextSteps:       "Next steps 3",
				},
			},
			expected: []template.Rule{
				{
					ID: "policy3",
					Properties: template.RuleProperties{
						Severity:  "LOW",
						NextSteps: "Next steps 3",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := constructRules(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				fmt.Printf("Expected %v, got %v\n", tc.expected, result)
			}
		})
	}
}

func TestConstructResults(t *testing.T) {
	testCases := []struct {
		name     string
		input    []template.Violation
		expected []template.Result
	}{
		{
			name:     "Empty Input",
			input:    []template.Violation{},
			expected: []template.Result{},
		},
		{
			name: "Missing Fields",
			input: []template.Violation{
				{
					PolicyID:      "policy1",
					AssetID:       "asset1",
					ViolatedAsset: template.AssetDetails{AssetType: "type1", Asset: "asset1"},
				},
			},
		},
		{
			name: "Multiple Violations",
			input: []template.Violation{
				{
					PolicyID:      "policy1",
					AssetID:       "asset1",
					NextSteps:     "next_steps1",
					ViolatedAsset: template.AssetDetails{AssetType: "type1", Asset: "asset1"},
				},
				{
					PolicyID:      "policy2",
					AssetID:       "asset2",
					NextSteps:     "next_steps2",
					ViolatedAsset: template.AssetDetails{AssetType: "type2", Asset: "asset2"},
				},
			},
			expected: []template.Result{
				{
					RuleID:  "policy1",
					Message: template.Message{Text: "Asset type: type1 has a violation, next steps: next_steps1"},
					Locations: []template.Location{
						{
							LogicalLocations: []template.LogicalLocation{
								{FullyQualifiedName: "asset1"},
							},
						},
					},
					Properties: template.ResultProperties{
						AssetID:   "asset1",
						Asset:     "asset1",
						AssetType: "type1",
					},
				},
				{
					RuleID:  "policy2",
					Message: template.Message{Text: "Asset type: type2 has a violation, next steps: next_steps2"},
					Locations: []template.Location{
						{
							LogicalLocations: []template.LogicalLocation{
								{FullyQualifiedName: "asset2"},
							},
						},
					},
					Properties: template.ResultProperties{
						AssetID:   "asset2",
						Asset:     "asset2",
						AssetType: "type2",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := constructResults(tc.input)

			if !reflect.DeepEqual(result, tc.expected) {
				fmt.Printf("Expected %v, got %v\n", tc.expected, result)
			}
		})
	}
}
