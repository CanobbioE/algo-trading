package monitor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ScanFilters defines criteria for filtering stocks.
type ScanFilters struct {
	MinConfidence    float64          `json:"min_confidence"`
	MinWeightedScore float64          `json:"min_weighted_score"`
	MaxRisk          RiskLevel        `json:"max_risk"`
	MinOpportunity   OpportunityLevel `json:"min_opportunity"`
	MinVolume        float64          `json:"min_volume"`
	MinMarketCap     float64          `json:"min_market_cap"`
	MaxMarketCap     float64          `json:"max_market_cap"`
	RequiredSignals  int              `json:"required_signals"`
}

// RiskLevel defines the probability of success for trading a product.
type RiskLevel int

const (
	// RiskLow represents a low probability of the trade being unsuccessful.
	RiskLow RiskLevel = iota
	// RiskMedium represents a medium probability of the trade being unsuccessful.
	RiskMedium
	// RiskHigh represents a high probability of the trade being unsuccessful.
	RiskHigh
)

// String returns the string representation of r.
func (r *RiskLevel) String() string {
	if r == nil {
		return "<nil>"
	}
	switch *r {
	case RiskLow:
		return "Low"
	case RiskMedium:
		return "Medium"
	case RiskHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// UnmarshalJSON implements a custom json.Unmarshaler.
func (r *RiskLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RiskLevel should be a string, got %s", data)
	}
	switch strings.ToUpper(s) {
	case "LOW":
		*r = RiskLow
	case "MEDIUM":
		*r = RiskMedium
	case "HIGH":
		*r = RiskHigh
	default:
		return fmt.Errorf("invalid RiskLevel: %q", s)
	}
	return nil
}

// OpportunityLevel defines how profitable a product is.
type OpportunityLevel int

const (
	// OpportunityLow represents a low margin for profit.
	OpportunityLow OpportunityLevel = iota
	// OpportunityMedium represents a medium margin for profit.
	OpportunityMedium
	// OpportunityHigh represents a high margin for profit.
	OpportunityHigh
)

// String returns the string representation of o.
func (o *OpportunityLevel) String() string {
	if o == nil {
		return "<nil>"
	}
	switch *o {
	case OpportunityLow:
		return "Low"
	case OpportunityMedium:
		return "Medium"
	case OpportunityHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// UnmarshalJSON implements a custom json.Unmarshaler.
func (o *OpportunityLevel) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("RiskLevel should be a string, got %s", data)
	}
	switch strings.ToUpper(s) {
	case "LOW":
		*o = OpportunityLow
	case "MEDIUM":
		*o = OpportunityMedium
	case "HIGH":
		*o = OpportunityHigh
	default:
		return fmt.Errorf("invalid RiskLevel: %q", s)
	}
	return nil
}
