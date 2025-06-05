package config

import (
	"encoding/json"
	"fmt"

	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
)

// Config is the project configuration.
type Config struct {
	Thresholds           *strategies.Thresholds       `json:"thresholds"`
	Filters              *monitor.ScanFilters         `json:"filters"`
	StockUniverse        []string                     `json:"stock_universe"`
	Strategies           []*strategies.StrategyWeight `json:"strategies"`
	LookBack             int                          `json:"lookback"`
	BollingerCoefficient float64                      `json:"bollinger_coefficient"`
}

type rawCfg struct {
	Thresholds           *strategies.Thresholds `json:"thresholds"`
	ScanFilters          *monitor.ScanFilters   `json:"scan_filters"`
	Strategies           []*rawStrategies       `json:"strategies"`
	StockUniverse        []string               `json:"stock_universe"`
	LookBack             int                    `json:"lookback"`
	BollingerCoefficient float64                `json:"bollinger_coefficient"`
}

type rawStrategies struct {
	Strategy string  `json:"strategy"`
	Weight   float64 `json:"weight"`
}

// UnmarshalJSON implements a custom json.Unmarshaler.
func (c *Config) UnmarshalJSON(data []byte) error {
	var raw rawCfg
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse json raw: %w", err)
	}
	c.Thresholds = raw.Thresholds
	c.Filters = raw.ScanFilters
	c.BollingerCoefficient = raw.BollingerCoefficient
	c.StockUniverse = raw.StockUniverse
	c.LookBack = raw.LookBack

	for _, str := range raw.Strategies {
		var s strategies.Strategy
		switch str.Strategy {
		case "BREAKOUT", "breakout":
			s = strategies.NewBreakoutStrategy(c.Thresholds)
		case "VWAP", "vwap":
			s = strategies.NewVWAPStrategy(c.LookBack)
		case "MEANREVERSION", "meanreversion", "meanReversion", "mean-reversion":
			s = strategies.NewMeanReversionStrategy(c.LookBack, c.Thresholds.Deviation)
		case "BOLLINGER", "bollinger":
			s = strategies.NewBollingerBandSqueezeStrategy(c.LookBack, c.BollingerCoefficient, c.Thresholds.Squeeze)
		}
		c.Strategies = append(c.Strategies, &strategies.StrategyWeight{
			Weight:   str.Weight,
			Strategy: s,
		})
	}

	return nil
}
