package strategies

import (
	"math"

	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type breakoutAnalysis struct {
	atr          float64
	resistance   float64
	support      float64
	latestVolume float64
	avgVolume    float64
}

// BreakoutStrategy implements a Breakout Strategy.
type BreakoutStrategy struct {
	t        *Thresholds
	analysis *breakoutAnalysis
}

// NewBreakoutStrategy creates a new BreakoutStrategy.
func NewBreakoutStrategy(t *Thresholds) *BreakoutStrategy {
	return &BreakoutStrategy{t: t}
}

// Execute scans the data and finds breakout signals based on high and low levels.
func (s *BreakoutStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.t.AtrPeriod+1 {
		return signals.NoOp
	}

	atr := calculateATR(data, s.t.AtrPeriod)
	var dynamicLookBack int
	switch {
	case atr < s.t.LowATRThreshold:
		dynamicLookBack = s.t.HighLookback
	case atr > s.t.HighATRThreshold:
		dynamicLookBack = s.t.LowLookback
	default:
		dynamicLookBack = (s.t.LowLookback + s.t.HighLookback) / 2
	}

	if len(data) < dynamicLookBack {
		return signals.NoOp
	}

	// Calculate breakout levels over dynamicLookBack
	var highestHigh, lowestLow float64
	for i := len(data) - dynamicLookBack; i < len(data); i++ {
		if i == len(data)-dynamicLookBack || data[i].High > highestHigh {
			highestHigh = data[i].High
		}
		if i == len(data)-dynamicLookBack || data[i].Low < lowestLow {
			lowestLow = data[i].Low
		}
	}

	// Calculate average volume over dynamicLookBack
	var totalVolume float64
	for i := len(data) - dynamicLookBack; i < len(data); i++ {
		totalVolume += data[i].Volume
	}
	avgVolume := totalVolume / float64(dynamicLookBack)

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

// calculateATR returns the Average True Range, which is the indicator that demonstrates market volatility.
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
