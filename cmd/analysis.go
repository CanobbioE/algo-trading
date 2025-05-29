package cmd

import (
	"fmt"
	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/signals"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
	"github.com/spf13/cobra"
	"strings"
)

type analysisScope struct {
	ticker    string
	timeFrame string
}

func init() {
	s := &analysisScope{}
	analysisCmd := &cobra.Command{
		Use:   "analyse",
		Short: "Analyse a single stock",
		Long:  "Analyse a single stock",
		RunE:  s.runE,
	}

	analysisCmd.Flags().StringVarP(&s.ticker, "ticker", "t", "", "Stock ticker to use")
	analysisCmd.Flags().StringVarP(&s.timeFrame, "timeframe", "f", "1d", "Time frame to use")

	utilities.Must(analysisCmd.MarkFlagRequired("ticker"))
	rootCmd.AddCommand(analysisCmd)
}

func (s *analysisScope) runE(_ *cobra.Command, _ []string) error {
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
	data, err := cli.GetOHLCV(s.ticker, &scraping.WithTimeframe{TimeFrame: scraping.TimeFrame(s.timeFrame)})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("no data for ticker %s", s.ticker)
	}

	m := make(map[signals.Operation]int, len(strats))

	for _, strat := range strats {
		operation := strat.Execute(data)
		m[operation]++
	}

	strategies.GenerateAnalysis(strategies.NewAnalysisInput(strats...))

	fmt.Println("Sentiment is:")
	for k, v := range m {
		var c color
		switch k {
		case signals.Buy:
			c = green
		case signals.Sell:
			c = blue
		case signals.Setup:
			c = yellow
		case signals.NoOp:
			c = white
		}
		printJSONColor(strings.ToUpper(string(k)), c)
		fmt.Printf(":\t%2.f%%\n", float32(v)/float32(len(strats))*100)
	}
	return nil
}
