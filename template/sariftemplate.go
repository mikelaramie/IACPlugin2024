package sarifTemplate

type SarifOutput struct {
	Version string `json:"version"`
	Schema  string `json:"$schema"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Note    string   `json:"note"`
	Tool    Tool     `json:"tool"`
	Results []Result `json:"results"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Driver struct {
	Name           string `json:"name"`
	Version        string `json:"version"`
	InformationURI string `json:"informationUri"`
	Rules          []Rule `json:"rules"`
}

type Rule struct {
	ID              string          `json:"id"`
	FullDescription FullDescription `json:"fullDescription"`
	Properties      RuleProperties  `json:"properties"`
}

type FullDescription struct {
	Text string `json:"text"`
}

type RuleProperties struct {
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

type Result struct {
	RuleID     string           `json:"ruleId"`
	Message    Message          `json:"message"`
	Locations  []Location       `json:"locations"`
	Properties ResultProperties `json:"properties"`
}

type Message struct {
	Text string `json:"text"`
}

type Location struct {
	LogicalLocations []LogicalLocation `json:"logicalLocation"`
}

type LogicalLocation struct {
	FullyQualifiedName string `json:"fullyQualifiedName"`
}

type ResultProperties struct {
	AssetID   string `json:"assetId"`
	AssetType string `json:"assetType"`
	Asset     string `json:"asset"`
}
