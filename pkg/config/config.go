package config

import (
	"encoding/json"
	"fmt"

	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
)

type Config struct {
	StockUniverse        []string                     `json:"stock_universe"`
	Strategies           []*strategies.StrategyWeight `json:"strategies"`
	Thresholds           *strategies.Thresholds       `json:"thresholds"`
	Filters              *monitor.ScanFilters         `json:"filters"`
	LookBack             int                          `json:"lookback"`
	BollingerCoefficient float64                      `json:"bollinger_coefficient"`
}

type rawCfg struct {
	Strategies []struct {
		Strategy string  `json:"strategy"`
		Weight   float64 `json:"weight"`
	} `json:"strategies"`
	StockUniverse        []string               `json:"stock_universe"`
	Thresholds           *strategies.Thresholds `json:"thresholds"`
	LookBack             int                    `json:"lookback"`
	BollingerCoefficient float64                `json:"bollinger_coefficient"`
	ScanFilters          *monitor.ScanFilters   `json:"scan_filters"`
}

func (c *Config) UnmarshalJSON(data []byte) error {
	var raw rawCfg
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse json raw: %w", err)
	}
	c.Thresholds = raw.Thresholds
	c.BollingerCoefficient = raw.BollingerCoefficient
	c.Filters = raw.ScanFilters
	c.StockUniverse = raw.StockUniverse

	for _, str := range raw.Strategies {
		var s strategies.Strategy
		switch str.Strategy {
		case "BREAKOUT":
			s = strategies.NewBreakouStrategy(c.Thresholds)
		case "VWAP":
			s = strategies.NewVWAPStrategy(c.LookBack)
		case "MEANREVERSION":
			s = strategies.NewMeanReversionStrategy(c.LookBack, c.Thresholds.Deviation)
		case "BOLLINGER":
			s = strategies.NewBollingerBandSqueezeStrategy(c.LookBack, c.BollingerCoefficient, c.Thresholds.Squeeze)
		}
		c.Strategies = append(c.Strategies, &strategies.StrategyWeight{
			Weight:   str.Weight,
			Strategy: s,
		})
	}

	return nil
}
