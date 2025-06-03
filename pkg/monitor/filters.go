package monitor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ScanFilters defines criteria for filtering stocks
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

func DefaultFilters() *ScanFilters {
	return &ScanFilters{
		MinConfidence:    0.6,
		MinWeightedScore: 1.0,
		MaxRisk:          RiskMedium,
		MinOpportunity:   OpportunityMedium,
		MinVolume:        100000,
		MinMarketCap:     1000000,     // $1M minimum
		MaxMarketCap:     50000000000, // $50B maximum
		RequiredSignals:  2,
	}
}

type RiskLevel int

const (
	RiskLow RiskLevel = iota
	RiskMedium
	RiskHigh
)

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

type OpportunityLevel int

const (
	OpportunityLow OpportunityLevel = iota
	OpportunityMedium
	OpportunityHigh
)

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
