package template

type SarifOutput struct {
	Version string `json:"version, omitempty"`
	Schema  string `json:"$schema, omitempty"`
	Runs    []Run  `json:"runs, omitempty"`
}

type Run struct {
	Note    string   `json:"note, omitempty"`
	Tool    Tool     `json:"tool, omitempty"`
	Results []Result `json:"results, omitempty"`
}

type Tool struct {
	Driver Driver `json:"driver, omitempty"`
}

type Driver struct {
	Name           string `json:"name, omitempty"`
	Version        string `json:"version, omitempty"`
	InformationURI string `json:"informationUri, omitempty"`
	Rules          []Rule `json:"rules, omitempty"`
}

type Rule struct {
	ID              string          `json:"id, omitempty"`
	FullDescription FullDescription `json:"fullDescription, omitempty"`
	Properties      RuleProperties  `json:"properties, omitempty"`
}

type FullDescription struct {
	Text string `json:"text, omitempty"`
}

type RuleProperties struct {
	Severity            string   `json:"severity, omitempty"`
	PolicyType          string   `json:"policyType, omitempty"`
	ComplianceStandard  []string `json:"complianceStandard, omitempty"`
	PolicySet           string   `json:"policySet, omitempty"`
	Posture             string   `json:"posture, omitempty"`
	PostureRevisionID   string   `json:"postureRevisionId, omitempty"`
	PostureDeploymentID string   `json:"postureDeploymentId, omitempty"`
	Constraints         string   `json:"constraints, omitempty"`
	NextSteps           string   `json:"nextSteps, omitempty"`
}

type Result struct {
	RuleID     string           `json:"ruleId, omitempty"`
	Message    Message          `json:"message, omitempty"`
	Locations  []Location       `json:"locations, omitempty"`
	Properties ResultProperties `json:"properties, omitempty"`
}

type Message struct {
	Text string `json:"text, omitempty"`
}

type Location struct {
	LogicalLocations []LogicalLocation `json:"logicalLocation, omitempty"`
}

type LogicalLocation struct {
	FullyQualifiedName string `json:"fullyQualifiedName, omitempty"`
}

type ResultProperties struct {
	AssetID   string `json:"assetId, omitempty"`
	AssetType string `json:"assetType, omitempty"`
	Asset     string `json:"asset, omitempty"`
}
