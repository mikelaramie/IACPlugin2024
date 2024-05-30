package utils

// IACScanReport is the json format passed as an input.
type IACScanReport struct {
	Response struct {
		IACValidationReport struct {
			Violations []struct {
				Severity string `json:"severity"`
			} `json:"violations"`
		} `json:"iacValidationReport"`
	} `json:"response"`
}
