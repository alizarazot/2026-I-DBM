package model

import (
	"encoding/json"
	"time"
)

// CFC stands for 'Customer Feedbacks and Complaints'.
type CFC struct {
	ID        string      `json:"id,omitempty"`
	Subject   string      `json:"subject"`
	Category  CFCCategory `json:"category"`
	UserEmail string      `json:"userEmail"`
	Details   string      `json:"details"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type CFCAnswer struct {
	ID        string    `json:"id,omitempty"`
	CFCID     string    `json:"cfcId"`
	UserEmail string    `json:"userId"`
	Answer    string    `json:"answer"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CFCCategory uint8

//go:generate go tool stringer -trimprefix CFCCategory -type CFCCategory

const (
	CFCCategoryInvalid CFCCategory = iota
	CFCCategoryRequest
	CFCCategoryComplaint
	CFCCategoryClaim
	CFCCategorySuggestion
)

func NewCFCCategory(canonical string) CFCCategory {
	switch canonical {
	case "request":
		return CFCCategoryRequest
	case "complaint":
		return CFCCategoryComplaint
	case "claim":
		return CFCCategoryClaim
	case "suggestion":
		return CFCCategorySuggestion
	default:
		return CFCCategoryInvalid
	}
}

func (c CFCCategory) CanonicalString() string {
	switch c {
	case CFCCategoryRequest:
		return "request"
	case CFCCategoryComplaint:
		return "complaint"
	case CFCCategoryClaim:
		return "claim"
	case CFCCategorySuggestion:
		return "suggestion"
	default:
		return "invalid"
	}
}

func (c CFCCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.CanonicalString())
}

func (c *CFCCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*c = NewCFCCategory(s)

	return nil
}
