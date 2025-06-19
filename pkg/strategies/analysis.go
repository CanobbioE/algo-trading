package strategies

import (
	"fmt"
	"math"

	"github.com/CanobbioE/stock-market-clients/api"

	"github.com/CanobbioE/algo-trading/pkg/printer"
	"github.com/CanobbioE/algo-trading/pkg/signals"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

// Thresholds collects all the shared limits for the various strategies.
type Thresholds struct {
	AtrPeriod         int     `json:"atr_period"`
	LowATRThreshold   float64 `json:"low_atr_threshold"`
	HighATRThreshold  float64 `json:"high_atr_threshold"`
	LowLookback       int     `json:"low_lookback"`
	HighLookback      int     `json:"high_lookback"`
	VolumeThreshold   float64 `json:"volume_threshold"`
	Deviation         float64 `json:"deviation"`
	Squeeze           float64 `json:"squeeze"`
	MinMomentumReturn float64 `json:"min_momentum_return"`
}

// Strategy is the common interface for all strategies.
type Strategy interface {
	// Execute the strategy with the given data and returns a suggestion in the form of a signals.Operation.
	Execute(data []*api.OHLCV) signals.Operation
}

// StrategyWeight encapsulate a strategy with a weight so that
// its results can be valued less or more compared to other strategies.
type StrategyWeight struct {
	Strategy Strategy
	Weight   float64
}

// Analysis the input parameters to perform a market analysis.
type Analysis struct {
	*vwapAnalysis
	*mrAnalysis
	*breakoutAnalysis
	*bbAnalysis
	*macdAnalysis
	*momentumAnalysis
	p printer.Printer
}

// NewAnalysisInput populate a new Analysis using the results from the given strategies.
func NewAnalysisInput(p printer.Printer, results ...*StrategyWeight) *Analysis {
	var out = &Analysis{
		p: p,
	}
	for _, r := range results {
		switch s := r.Strategy.(type) {
		case *VWAPStrategy:
			out.vwapAnalysis = utilities.DefaultPointer(s.analysis)
		case *MeanReversionStrategy:
			out.mrAnalysis = utilities.DefaultPointer(s.analysis)
		case *BreakoutStrategy:
			out.breakoutAnalysis = utilities.DefaultPointer(s.analysis)
		case *BollingerBandSqueezeStrategy:
			out.bbAnalysis = utilities.DefaultPointer(s.analysis)
		case *MACDStrategy:
			out.macdAnalysis = utilities.DefaultPointer(s.analysis)
		case *MomentumStrategy:
			out.momentumAnalysis = utilities.DefaultPointer(s.analysis)
		}
	}

	return out
}

// GenerateAnalysis pretty prints a human-readable report from the given Analysis.
func (in *Analysis) GenerateAnalysis() {
	in.p.Println("Key Technical Signals:")
	in.p.Println("======================")

	// Close vs VWAP/SMA
	in.p.Printf("- Latest Close (%.3f) is %s VWAP (%.3f)\n\t-> %s\n",
		in.closePrice, compare(in.closePrice, in.vwap), in.vwap, closeOverVWAP(in.closePrice, in.vwap))
	in.p.Printf("- Latest Close (%.3f) is %s SMA (%.3f)\n\t-> %s\n",
		in.closePrice, compare(in.closePrice, in.sma), in.sma, closeOverSMA(in.closePrice, in.sma))

	// Deviation
	in.p.Printf("- Deviation from SMA: %.2f%%\n\t-> %s\n", in.deviation, trendStrength(in.deviation))

	// Resistance/Support
	in.p.Printf("- Resistance Level: %.3f, Support Level: %.3f\n", in.resistance, in.support)
	in.p.Printf("- Price is currently %s breakout range\n",
		breakoutStatus(in.closePrice, in.support, in.resistance))
	if in.latestVolume > in.avgVolume {
		in.p.PrintColored(printer.Green, "- Volume is above average\n\t-> Validates potential breakout\n")
	} else {
		in.p.PrintColored(printer.Yellow, "- Volume is below average\n\t-> Watch for confirmation before acting")
	}

	// RSI
	in.p.Printf("- Relative Strength: %.2f\n\t-> %s\n", in.rsi, rsiStatus(in.rsi))

	// Bollinger Bands
	in.p.Printf("- Bollinger Band SMA: %.3f, Upper: %.3f, Lower: %.3f, Width: %.3f\n\t-> %s\n",
		in.bbSMA, in.upper, in.lower, in.width, bollingerSuggestion(in.closePrice, in.upper, in.lower))

	// ATR
	in.p.Printf("- ATR: %.4f\n\t-> Volatility is %s\n", in.atr, volatilityStatus(in.atr, in.width))

	// MACD
	in.p.Println(macdSuggestion(in.prevDelta, in.delta, in.triggerDistance))

	// Momentum
	in.p.Printf("- Found a %.2f%% price change in the last %d periods\n\t -> %s\n",
		in.change*100, in.period, momentumStatus(in.change))
	in.p.Println("======================")
}

func compare(a, b float64) string {
	if a > b {
		return printer.WrapInColor("above", printer.Green)
	} else if a < b {
		return printer.WrapInColor("below", printer.Red)
	}
	return "equal to"
}

func trendStrength(deviation float64) string {
	if deviation < -2 {
		return printer.WrapInColor("trend strength: slight weakness", printer.Red)
	} else if deviation > 2 {
		return printer.WrapInColor("trend strength: slight strength", printer.Green)
	}
	return "trend strength: near neutral"
}

func breakoutStatus(price, support, resistance float64) string {
	if price > resistance {
		return printer.WrapInColor("above resistance (possible breakout)", printer.Green)
	} else if price < support {
		return printer.WrapInColor("below support (possible breakdown)", printer.Red)
	}
	return "inside"
}

func rsiStatus(rsi float64) string {
	switch {
	case rsi < 30:
		return printer.WrapInColor("Oversold (reversal opportunity)", printer.Green)
	case rsi > 70:
		return printer.WrapInColor("Overbought (sell and take profit)", printer.Blue)
	default:
		return "Neutral (sideway market)"
	}
}

func volatilityStatus(atr, bbWidth float64) string {
	if bbWidth < 0.035 && atr < 0.02 {
		return "low"
	}
	return "moderate or high"
}

func closeOverVWAP(cls, vwap float64) string {
	if cls > vwap {
		return printer.WrapInColor("Suggests buyers are in control (BULLISH)", printer.Green)
	}
	if cls < vwap {
		return printer.WrapInColor("Suggests selling pressure, possibly distribution (BEARISH)", printer.Red)
	}

	return "No intraday edge for bulls or bears"
}

func closeOverSMA(cls, sma float64) string {
	if cls > sma {
		return printer.WrapInColor("Indicates uptrend continuation or strength (support in rising trend)", printer.Green)
	}
	if cls < sma {
		return printer.WrapInColor("Suggests weakness/downtrend - "+
			"Can be a bearish sign, especially if price breaks down from a key moving average", printer.Red)
	}
	return "Trend indecision, can act as a magnet"
}

func bollingerSuggestion(currentPrice, upper, lower float64) string {
	const threshold = 0.01 // 1% threshold to consider "near" the band
	switch {
	case currentPrice >= upper, math.Abs(currentPrice-upper) <= upper*threshold:
		return printer.WrapInColor("Consider selling (price near or above upper band)", printer.Blue)
	case currentPrice <= lower, math.Abs(currentPrice-lower) <= lower*threshold:
		return printer.WrapInColor("Consider buying (price near or below lower band)", printer.Green)
	default:
		return "Hold (price within bands)"
	}
}

func macdSuggestion(prevDelta, delta, triggerDistance float64) string {
	var (
		crossover  string
		suggestion string
	)
	switch {
	case prevDelta < 0 && delta > 0:
		crossover = "MACD crossover: " + printer.WrapInColor("Bullish", printer.Green) +
			" signal detected (MACD crossed above signal line)."
		suggestion = printer.WrapInColor("Consider buying", printer.Green)
	case prevDelta > 0 && delta < 0:
		crossover = "MACD crossover: " + printer.WrapInColor("Bearish", printer.Red) +
			" signal detected (MACD crossed below signal line)."
		suggestion = printer.WrapInColor("Consider selling", printer.Blue)
	case delta > 0:
		crossover = "MACD shows " + printer.WrapInColor("bullish", printer.Green) + " momentum."
		suggestion = printer.WrapInColor("Consider buying", printer.Green) + " (if confirmed by other indicators)"
	case delta < 0:
		crossover = "MACD shows " + printer.WrapInColor("bearish", printer.Red) + " momentum."
		suggestion = printer.WrapInColor("Consider selling", printer.Blue) + " (if confirmed by other indicators)"
	default:
		crossover = "MACD is neutral, showing no significant momentum or crossover."
		suggestion = "Hold"
	}

	if delta == 0 {
		return crossover
	}

	str := "weak"
	dir := "below"
	if math.Abs(triggerDistance) > 0.5 {
		str = "strong"
	}
	if triggerDistance > 0 {
		dir = "above"
	}

	var bias string
	switch {
	case str == "strong" && dir == "above":
		str = printer.WrapInColor(str, printer.Green)
		bias = printer.WrapInColor("buy", printer.Green)
	case str == "strong":
		str = printer.WrapInColor(str, printer.Blue)
		bias = printer.WrapInColor("sell", printer.Blue)
	case dir == "above":
		bias = "hold"
	default:
		str = printer.WrapInColor(str, printer.Yellow)
		bias = printer.WrapInColor("caution", printer.Yellow)
	}

	return fmt.Sprintf("- %s\n\t-> Momentum strength is %s with MACD %s the zero line (%s bias).\n\t-> %s",
		crossover, str, dir, bias, suggestion)
}

func momentumStatus(change float64) string {
	var direction string
	switch {
	case change > 0.05:
		direction = printer.WrapInColor("strong upward", printer.Green)
	case change > 0.02:
		direction = printer.WrapInColor("moderate upward", printer.Green)
	case change > 0:
		direction = printer.WrapInColor("slight upward", printer.Green)
	case change < -0.05:
		direction = printer.WrapInColor("strong downward", printer.Red)
	case change < -0.02:
		direction = printer.WrapInColor("moderate downward", printer.Red)
	case change < 0:
		direction = printer.WrapInColor("slight downward", printer.Red)
	default:
		direction = "no significant"
	}
	return "Indicates a " + direction + " momentum."
}
