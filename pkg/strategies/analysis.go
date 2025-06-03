package strategies

import (
	"fmt"
	"math"

	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/signals"
)

type Strategy interface {
	Execute(data []*api.OHLCV) signals.Operation
}

type StrategyWeight struct {
	Weight   float64
	Strategy Strategy
}

type AnalysisInput struct {
	*vwapAnalysis
	*mrAnalysis
	*breakoutAnalysis
	*bbAnalysis
}

func NewAnalysisInput(results ...Strategy) *AnalysisInput {
	var out = &AnalysisInput{}
	for _, r := range results {
		switch s := r.(type) {
		case *VWAPStrategy:
			out.vwapAnalysis = s.analysis
		case *MeanReversionStrategy:
			out.mrAnalysis = s.analysis
		case *BreakouStrategy:
			out.breakoutAnalysis = s.analysis
		case *BollingerBandSqueezeStrategy:
			out.bbAnalysis = s.analysis
		}
	}

	return out
}

func GenerateAnalysis(in *AnalysisInput) {
	fmt.Println("Key Technical Signals:")
	fmt.Println("======================")

	// Close vs VWAP/SMA
	fmt.Printf("- Latest Close (%.3f) is %s VWAP (%.3f)\n\t-> %s\n", in.closePrice, compare(in.closePrice, in.vwap), in.vwap, closeOverVWAP(in.closePrice, in.vwap))
	fmt.Printf("- Latest Close (%.3f) is %s SMA (%.3f)\n\t-> %s\n", in.closePrice, compare(in.closePrice, in.sma), in.sma, closeOverSMA(in.closePrice, in.sma))

	// Deviation
	fmt.Printf("- Deviation from SMA: %.2f%%\n\t-> %s\n", in.deviation, trendStrength(in.deviation))

	// Resistance/Support
	fmt.Printf("- Resistance Level: %.3f, Support Level: %.3f\n", in.resistance, in.support)
	fmt.Printf("- Price is currently %s breakout range\n", breakoutStatus(in.closePrice, in.support, in.resistance))
	if in.latestVolume > in.avgVolume {
		fmt.Println("- Volume is above average\n\t-> Validates potential breakout")
	} else {
		fmt.Println("- Volume is below average\n\t-> Watch for confirmation before acting")
	}

	// RSI
	fmt.Printf("- Relative Strength: %.2f\n\t-> %s\n", in.rsi, rsiStatus(in.rsi))

	// Bollinger Bands
	fmt.Printf("- Bollinger Band SMA: %.3f, Upper: %.3f, Lower: %.3f, Width: %.3f\n\t-> %s\n", in.bbSMA, in.upper, in.lower, in.width, bollingerSuggestion(in.closePrice, in.upper, in.lower))

	// ATR
	fmt.Printf("- ATR: %.4f\n\t-> Volatility is %s\n", in.atr, volatilityStatus(in.atr, in.width))
	fmt.Println("======================")

}

func compare(a, b float64) string {
	if a > b {
		return "above"
	} else if a < b {
		return "below"
	}
	return "equal to"
}

func trendStrength(deviation float64) string {
	if deviation < -2 {
		return "trend strength: slight weakness"
	} else if deviation > 2 {
		return "trend strength: slight strength"
	}
	return "trend strength: near neutral"
}

func breakoutStatus(price, support, resistance float64) string {
	if price > resistance {
		return "above resistance (possible breakout)"
	} else if price < support {
		return "below support (possible breakdown)"
	}
	return "inside"
}

func rsiStatus(rsi float64) string {
	switch {
	case rsi < 30:
		return "Oversold (reversal opportunity)"
	case rsi > 70:
		return "Overbought (sell and take profit)"
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

func closeOverVWAP(close, vwap float64) string {
	if close > vwap {
		return "Suggests buyers are in control (BULLISH)"

	}
	if close < vwap {
		return "Suggests selling pressure, possibly distribution (BEARISH)"
	}

	return "No intraday edge for bulls or bears"
}

func closeOverSMA(close, sma float64) string {
	if close > sma {
		return "Indicates uptrend continuation or strength (support in rising trend)"
	}
	if close < sma {
		return "Suggests weakness/downtrend - Can be a bearish sign, especially if price breaks down from a key moving average"
	}
	return "Trend indecision, can act as a magnet"
}

func bollingerSuggestion(currentPrice, upper, lower float64) string {
	const threshold = 0.01 // 1% threshold to consider "near" the band

	if currentPrice >= upper || math.Abs(currentPrice-upper) <= upper*threshold {
		return "Consider selling (price near or above upper band)"
	} else if currentPrice <= lower || math.Abs(currentPrice-lower) <= lower*threshold {
		return "Consider buying (price near or below lower band)"
	} else {
		return "Hold (price within bands)"
	}
}
