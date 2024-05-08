package sarif

type Violation struct {
	AssetID        string `json:"assetId"`
	NextSteps      string `json:"nextSteps"`
	PolicyID       string `json:"policyId"`
	Severity       string `json:"severity"`
	ViolatedAsset  Asset  `json:"violatedAsset"`
	ViolatedPolicy Policy `json:"violatedPolicy"`
}

type Asset struct {
	Asset     string `json:"asset"`
	AssetType string `json:"assetType"`
}

type Policy struct {
	Constraint     string `json:"constraint"`
	ConstraintType string `json:"constraintType"`
}

type SarifOutput struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []SarifRun `json:"runs"`
}

type SarifRun struct {
	Tool    Tool          `json:"tool"`
	Results []SarifResult `json:"results"`
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

type SarifResult struct {
	RuleID     string         `json:"ruleId"`
	Message    Message        `json:"message"`
	Locations  []Location     `json:"locations"`
	Properties PropertyResult `json:"properties"`
}

type PropertyResult struct {
	AssetID   string `json:"assetId"`
	AssetType string `json:"assetType"`
	Asset     string `json:"asset"`
}

type Message struct {
	Text string `json:"text"`
}

type Rule struct {
	ID              string          `json:"id"`
	FullDescription FullDescription `json:"fullDescription"`
	Properties      PropertyRule    `json:"properties"`
}

type PropertyRule struct {
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

type FullDescription struct {
	Text string `json:"text"`
}

type Location struct {
	LogicalLocation []LogicalLocation `json:"logicalLocation"`
}

type LogicalLocation struct {
	FullyQualifiedName []string `json:"fullyQualifiedName"`
}

type physicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
}

type ArtifactLocation struct {
	URI string `json:"uri"`
}
