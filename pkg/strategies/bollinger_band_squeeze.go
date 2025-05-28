package strategies

import (
	"fmt"
	"math"

	"github.com/canobbioe/algo-trading/pkg/api"
	"github.com/canobbioe/algo-trading/pkg/signals"
)

type bbAnalysis struct {
	bbSMA float64
	upper float64
	lower float64
	width float64
}
type BollingerBandSqueezeStrategy struct {
	period           int
	k                float64
	squeezeThreshold float64
	analysis         *bbAnalysis
}

func NewBollingerBandSqueezeStrategy(period int, k, squeezeThreshold float64) *BollingerBandSqueezeStrategy {
	return &BollingerBandSqueezeStrategy{
		period:           period,
		k:                k,
		squeezeThreshold: squeezeThreshold,
	}
}

func calculateSMA(data []*api.OHLCV) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, bar := range data {
		sum += bar.Close
	}
	return sum / float64(len(data))
}

func calculateStdDev(data []*api.OHLCV, sma float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sumSquares := 0.0
	for _, bar := range data {
		diff := bar.Close - sma
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(data))
	return math.Sqrt(variance)
}

// Execute Bollinger Band Squeeze Strategy.
// This strategy detects when volatility is low (bands squeeze), and waits for a breakout as the bands expand.
func (s *BollingerBandSqueezeStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.period {
		return signals.NoOp
	}

	recentData := data[len(data)-s.period:]
	sma := calculateSMA(recentData)
	stdDev := calculateStdDev(recentData, sma)

	upperBand := sma + s.k*stdDev
	lowerBand := sma - s.k*stdDev
	bandWidth := upperBand - lowerBand

	s.analysis = &bbAnalysis{
		bbSMA: sma,
		upper: upperBand,
		lower: lowerBand,
		width: bandWidth,
	}

	// Determine if the bandwidth is below threshold (squeeze)
	if bandWidth < s.squeezeThreshold {
		fmt.Println("Bollinger Band Squeeze detected -> Potential breakout setup")
		// Check for breakout
		latest := data[len(data)-1]
		if latest.Close > upperBand {
			return signals.Buy
		} else if latest.Close < lowerBand {
			return signals.Sell
		}
		return signals.Setup
	}

	return signals.NoOp
}
