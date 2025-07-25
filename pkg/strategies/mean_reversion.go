package strategies

import (
	"math"

	"github.com/CanobbioE/stock-market-clients/api"

	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type mrAnalysis struct {
	rsi       float64
	deviation float64
	sma       float64
}

// MeanReversionStrategy implements a Mean Reversion Strategy.
type MeanReversionStrategy struct {
	analysis           *mrAnalysis
	lookBack           int
	rsiPeriod          int
	deviationThreshold float64
}

// NewMeanReversionStrategy creates a new MeanReversionStrategy.
func NewMeanReversionStrategy(lookBack int, deviationThreshold float64) *MeanReversionStrategy {
	return &MeanReversionStrategy{
		lookBack:           lookBack,
		deviationThreshold: deviationThreshold,
	}
}

// Execute mean reversion with RSI confirmation.
func (s *MeanReversionStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.lookBack || len(data) <= s.rsiPeriod {
		return signals.NoOp
	}

	// Calculate SMA
	start := len(data) - s.lookBack
	sum := 0.0
	for i := start; i < len(data); i++ {
		sum += data[i].Close
	}
	sma := sum / float64(s.lookBack)

	latest := data[len(data)-1]
	deviation := (latest.Close - sma) / sma

	// Calculate RSI
	rsi := s.calculateRSI(data, s.rsiPeriod)

	s.analysis = &mrAnalysis{
		rsi:       rsi,
		deviation: deviation,
		sma:       sma,
	}

	if deviation < -s.deviationThreshold && rsi < 30 {
		return signals.Buy // Oversold and below mean
	} else if deviation > s.deviationThreshold && rsi > 70 {
		return signals.Sell // Overbought and above mean
	}

	return signals.NoOp
}

// calculateRSI computes RSI over a s.lookBack period.
// RSI measures momentum by comparing average gains and losses over a period (e.g., 14 bars).
// RSI ranges between 0–100:
// - Above 70: Overbought -> Possible sell signal.
// - Below 30: Oversold -> Possible buy signal.
func (*MeanReversionStrategy) calculateRSI(data []*api.OHLCV, period int) float64 {
	if len(data) <= period {
		return 50 // Neutral RSI
	}

	var gains, losses float64
	for i := len(data) - period; i < len(data); i++ {
		change := data[i].Close - data[i-1].Close
		if change > 0 {
			gains += change
		} else {
			losses += math.Abs(change)
		}
	}

	if gains+losses == 0 {
		return 50
	}

	rs := gains / losses
	rsi := 100 - (100 / (1 + rs))
	return rsi
}
