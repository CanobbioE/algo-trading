package strategies

import (
	"github.com/CanobbioE/stock-market-clients/api"

	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type vwapAnalysis struct {
	closePrice float64
	vwap       float64
}

// VWAPStrategy implements a Volume Weighted Average Price Strategy.
type VWAPStrategy struct {
	analysis *vwapAnalysis
	lookBack int
}

// NewVWAPStrategy creates a new VWAPStrategy.
func NewVWAPStrategy(lookBack int) *VWAPStrategy {
	return &VWAPStrategy{lookBack: lookBack}
}

// Execute calculates VWAP and returns buy/sell/none decision.
func (s *VWAPStrategy) Execute(data []*api.OHLCV) signals.Operation {
	if len(data) == 0 || s.lookBack <= 0 {
		return signals.NoOp
	}

	start := len(data) - s.lookBack
	if start < 0 {
		start = 0 // Use available data if lookback exceeds size
	}

	var cumulativePV, cumulativeVolume float64
	for i := start; i < len(data); i++ {
		bar := data[i]
		typicalPrice := (bar.High + bar.Low + bar.Close) / 3
		pv := typicalPrice * bar.Volume
		cumulativePV += pv
		cumulativeVolume += bar.Volume
	}

	if cumulativeVolume == 0 {
		return signals.NoOp
	}

	vwap := cumulativePV / cumulativeVolume
	latest := data[len(data)-1]

	s.analysis = &vwapAnalysis{
		closePrice: latest.Close,
		vwap:       vwap,
	}
	if latest.Close > vwap {
		return signals.Buy
	} else if latest.Close < vwap {
		return signals.Sell
	}
	return signals.NoOp
}
