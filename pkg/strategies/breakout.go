package strategies

import (
	"fmt"
	"math"

	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type Thresholds struct {
	AtrPeriod        int     `json:"atr_period"`
	LowATRThreshold  float64 `json:"low_atr_threshold"`
	HighATRThreshold float64 `json:"high_atr_threshold"`
	LowLookback      int     `json:"low_lookback"`
	HighLookback     int     `json:"high_lookback"`
	VolumeThreshold  float64 `json:"volume_threshold"`
	Deviation        float64 `json:"deviation"`
	Squeeze          float64 `json:"squeeze"`
}

type breakoutAnalysis struct {
	atr          float64
	resistance   float64
	support      float64
	latestVolume float64
	avgVolume    float64
}

type BreakouStrategy struct {
	t        *Thresholds
	analysis *breakoutAnalysis
}

func NewBreakouStrategy(t *Thresholds) *BreakouStrategy {
	return &BreakouStrategy{t: t}
}

// Execute scans the data and finds breakout signals based on high and low levels
func (s *BreakouStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.t.AtrPeriod+1 {
		return signals.NoOp
	}

	// Calculate ATR
	atr := calculateATR(data, s.t.AtrPeriod)
	//fmt.Printf("Current ATR: %.4f\n", atr)

	// Determine dynamic lookback
	var lookback int
	if atr < s.t.LowATRThreshold {
		lookback = s.t.HighLookback
	} else if atr > s.t.HighATRThreshold {
		lookback = s.t.LowLookback
	} else {
		lookback = (s.t.LowLookback + s.t.HighLookback) / 2
	}
	fmt.Printf("Dynamic Lookback Period: %d\n", lookback)

	if len(data) < lookback {
		return signals.NoOp
	}

	// Calculate breakout levels over lookback
	var highestHigh, lowestLow float64
	for i := len(data) - lookback; i < len(data); i++ {
		if i == len(data)-lookback || data[i].High > highestHigh {
			highestHigh = data[i].High
		}
		if i == len(data)-lookback || data[i].Low < lowestLow {
			lowestLow = data[i].Low
		}
	}

	// Calculate average volume over lookback
	var totalVolume float64
	for i := len(data) - lookback; i < len(data); i++ {
		totalVolume += data[i].Volume
	}
	avgVolume := totalVolume / float64(lookback)
	//fmt.Printf("Average volume for confirmation: %.0f\n", avgVolume)

	latest := data[len(data)-1]
	s.analysis = &breakoutAnalysis{
		atr:          atr,
		resistance:   highestHigh,
		support:      lowestLow,
		latestVolume: latest.Volume,
		avgVolume:    avgVolume,
	}
	// Check most recent bar for breakout
	if latest.Close > highestHigh && latest.Volume > avgVolume*s.t.VolumeThreshold {
		return signals.Buy
	} else if latest.Close < lowestLow && latest.Volume > avgVolume*s.t.VolumeThreshold {
		return signals.Sell
	}

	return signals.NoOp
}

func calculateATR(data []*api.OHLCV, period int) float64 {
	if len(data) < period+1 {
		return 0
	}
	var sumATR float64
	for i := len(data) - period; i < len(data); i++ {
		tr := math.Max(data[i].High-data[i].Low, math.Max(
			math.Abs(data[i].High-data[i-1].Close),
			math.Abs(data[i].Low-data[i-1].Close),
		))
		sumATR += tr
	}
	return sumATR / float64(period)
}
