package strategies

import (
	"github.com/CanobbioE/stock-market-clients/api"

	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type momentumAnalysis struct {
	change float64
	period int
}

// MomentumStrategy implements a Momentum Strategy.
type MomentumStrategy struct {
	t        *Thresholds
	analysis *momentumAnalysis
	lookBack int
}

// NewMomentumStrategy creates a new MomentumStrategy.
func NewMomentumStrategy(lookBack int, t *Thresholds) *MomentumStrategy {
	return &MomentumStrategy{
		lookBack: lookBack,
		t:        t,
	}
}

// Execute calculates the percentage change over a lookback period.
// If the return is above the minReturn threshold, it signals a BUY.
func (s *MomentumStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) < s.lookBack+1 {
		return signals.NoOp
	}
	start := data[len(data)-s.lookBack-1]
	end := data[len(data)-1]

	s.analysis = &momentumAnalysis{
		change: (end.Close - start.Close) / start.Close,
		period: s.lookBack,
	}

	switch {
	case s.analysis.change > s.t.MinMomentumReturn:
		return signals.Buy
	case s.analysis.change < -s.t.MinMomentumReturn:
		return signals.Sell
	default:
		return signals.NoOp
	}
}
