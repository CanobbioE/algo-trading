package strategies

import (
	"github.com/CanobbioE/stock-market-clients/api"

	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type macdAnalysis struct {
	prevDelta       float64
	delta           float64
	triggerDistance float64
}

// MACDParams defines input parameters for the MACDStrategy.
type MACDParams struct {
	FastPeriod      int     `json:"fast_period"`
	SlowPeriod      int     `json:"slow_period"`
	SignalPeriod    int     `json:"signal_period"`
	TriggerDistance float64 `json:"trigger_distance"`
}

// MACDStrategy implements a Moving Average Convergence Divergence Strategy that
// suggests to BUY if MACD crosses above signal line, SELL if below.
type MACDStrategy struct {
	analysis *macdAnalysis
	in       *MACDParams
}

// calculateEMA calculates the Exponential Moving Average for a given period.
func calculateEMA(data []*api.OHLCV, period int) []float64 {
	alpha := 2.0 / float64(period+1)
	ema := make([]float64, len(data))
	ema[0] = data[0].Close // seed with first value
	for i := 1; i < len(data); i++ {
		ema[i] = alpha*data[i].Close + (1-alpha)*ema[i-1]
	}
	return ema
}

// NewMACDStrategy creates a new MACDStrategy.
func NewMACDStrategy(in *MACDParams) *MACDStrategy {
	return &MACDStrategy{
		in: in,
	}
}

// Execute identifies if the current price is above or below the moving average.
func (s *MACDStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.in.SlowPeriod+s.in.SignalPeriod || len(data) < 2 {
		return signals.NoOp
	}

	fastEMA := calculateEMA(data, s.in.FastPeriod)
	slowEMA := calculateEMA(data, s.in.SlowPeriod)

	// Compute MACD line
	macdLine := make([]*api.OHLCV, len(data))
	for i := range data {
		macdLine[i] = &api.OHLCV{
			Close: fastEMA[i] - slowEMA[i],
		}
	}

	// Compute Signal line (EMA of MACD)
	signalLine := calculateEMA(macdLine, s.in.SignalPeriod)

	// Look at the last two values to detect crossover
	n := len(data)
	prevMACD := macdLine[n-2]
	currMACD := macdLine[n-1]
	prevSignal := signalLine[n-2]
	currSignal := signalLine[n-1]

	s.analysis = &macdAnalysis{
		prevDelta:       prevMACD.Close - prevSignal,
		delta:           currMACD.Close - currSignal,
		triggerDistance: s.in.TriggerDistance,
	}

	switch {
	case s.analysis.prevDelta < 0 && s.analysis.delta > s.in.TriggerDistance:
		return signals.Buy
	case s.analysis.prevDelta > 0 && s.analysis.delta < -s.in.TriggerDistance:
		return signals.Sell
	default:
		return signals.NoOp
	}
}
