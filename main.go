package main

import (
	"fmt"
	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/signals"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
)

func main() {
	lookBack := 3
	t := &strategies.Thresholds{
		AtrPeriod:        3,
		LowATRThreshold:  0.05,
		HighATRThreshold: 0.2,
		LowLookback:      10,
		HighLookback:     3,
		VolumeThreshold:  1.2,
	}

	// create a bunch of strategies:
	strats := []strategies.Strategy{
		strategies.NewBreakouStrategy(t),
		strategies.NewVWAPStrategy(lookBack),
		strategies.NewMeanReversionStrategy(lookBack, 0.02),
		strategies.NewBollingerBandSqueezeStrategy(lookBack, 2.0, 0.1),
	}

	cli := scraping.NewClient()
	data, err := cli.GetOHLCV("QQQS.ETF", &scraping.WithTimeframe{TimeFrame: scraping.Daily})
	if err != nil {
		panic(err)
	}

	m := make(map[signals.Operation]int, len(strats))

	for _, strat := range strats {
		operation := strat.Execute(data)
		m[operation]++
	}

	strategies.GenerateAnalysis(strategies.NewAnalysisInput(strats...))

	fmt.Println("sentiment is:")
	for k, v := range m {
		fmt.Printf("%s: %d\n", k, v)
	}
}
