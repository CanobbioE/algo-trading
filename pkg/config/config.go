package config

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
)

type (
	// Config is the project configuration.
	Config struct {
		Thresholds           *strategies.Thresholds       `json:"thresholds"`
		MACDParams           *strategies.MACDParams       `json:"macd_params"`
		Filters              *monitor.ScanFilters         `json:"filters"`
		StockUniverse        []string                     `json:"stock_universe"`
		Strategies           []*strategies.StrategyWeight `json:"strategies"`
		LookBack             int                          `json:"lookback"`
		MomentumLookBack     int                          `json:"momentum_look_back"`
		BollingerCoefficient float64                      `json:"bollinger_coefficient"`
	}
	rawCfg struct {
		Thresholds           *strategies.Thresholds `json:"thresholds"`
		MACDParams           *strategies.MACDParams `json:"macd_params"`
		ScanFilters          *monitor.ScanFilters   `json:"scan_filters"`
		Strategies           []*rawStrategies       `json:"strategies"`
		StockUniverse        []string               `json:"stock_universe"`
		LookBack             int                    `json:"lookback"`
		MomentumLookBack     int                    `json:"momentum_lookback"`
		BollingerCoefficient float64                `json:"bollinger_coefficient"`
	}
	rawStrategies struct {
		Strategy string  `json:"strategy"`
		Weight   float64 `json:"weight"`
	}
)

// UnmarshalJSON implements a custom json.Unmarshaler.
func (c *Config) UnmarshalJSON(data []byte) error {
	var raw rawCfg
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to parse json raw: %w", err)
	}
	c.Thresholds = raw.Thresholds
	if c.Thresholds == nil {
		return errors.New("no thresholds specified")
	}
	c.Filters = raw.ScanFilters
	c.BollingerCoefficient = raw.BollingerCoefficient
	c.StockUniverse = raw.StockUniverse
	c.LookBack = raw.LookBack
	c.MACDParams = raw.MACDParams
	c.MomentumLookBack = raw.MomentumLookBack

	c.Strategies = make([]*strategies.StrategyWeight, 0, len(raw.Strategies))
	if len(raw.Strategies) == 0 {
		return errors.New("at least one strategy must be specified")
	}

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
		case "MACD", "macd":
			s = strategies.NewMACDStrategy(c.MACDParams)
		case "MOMENTUM", "momentum":
			s = strategies.NewMomentumStrategy(c.MomentumLookBack, c.Thresholds)
		default:
			return fmt.Errorf("unknown strategy %s", str.Strategy)
		}
		c.Strategies = append(c.Strategies, &strategies.StrategyWeight{
			Weight:   str.Weight,
			Strategy: s,
		})
	}

	return nil
}
