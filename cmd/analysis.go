package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/CanobbioE/stock-market-clients/api"
	"github.com/CanobbioE/stock-market-clients/carnost"
	"github.com/spf13/cobra"

	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/printer"
	"github.com/CanobbioE/algo-trading/pkg/signals"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

type analysisScope struct {
	p           printer.Printer
	cfg         *config.Config
	ticker      string
	timeFrame   string
	mode        string
	cfgFile     string
	refreshRate time.Duration
	lifespan    time.Duration
}

func init() {
	s := &analysisScope{
		p: &printer.Standard{},
	}
	analysisCmd := &cobra.Command{
		Use:     "analyse",
		Short:   "Analyse a single stock",
		Long:    "Analyse a single stock",
		RunE:    s.runE,
		PreRunE: s.preRunE,
	}

	analysisCmd.Flags().StringVarP(&s.cfgFile, "config", "c", "", "Path to config file")
	analysisCmd.Flags().StringVarP(&s.ticker, "ticker", "t", "", "Stock ticker to use")
	analysisCmd.Flags().StringVarP(&s.timeFrame, "timeframe", "f", "1d", "Time frame to use")
	analysisCmd.Flags().StringVarP(&s.mode, "mode", "m", "onetime", "How the command will run: continue or onetime")
	//nolint:lll
	analysisCmd.Flags().DurationVarP(&s.refreshRate, "refresh", "r", 10*time.Minute, "Analyses refresh rate in continuos mode")
	analysisCmd.Flags().DurationVarP(&s.lifespan, "life", "l", 1*time.Hour, "How long the continuous mode should run for")

	utilities.Must(analysisCmd.MarkFlagRequired("ticker"))
	utilities.Must(analysisCmd.MarkFlagRequired("config"))
	rootCmd.AddCommand(analysisCmd)
}

func (s *analysisScope) preRunE(_ *cobra.Command, _ []string) error {
	file, err := os.Open(s.cfgFile)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	var cfg config.Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}

	s.cfg = &cfg
	return nil
}

func (s *analysisScope) runE(cmd *cobra.Command, _ []string) error {
	cli := carnost.NewClient()
	ctx := cmd.Context()

	switch s.mode {
	case "onetime", "once", "one-time":
		return s.analyse(ctx, cli)
	case "continuous", "continue", "cont":
		if s.refreshRate < 5*time.Second {
			s.p.PrintColored(printer.Yellow, "Refresh rate of %s is too low, defaulting to 5s\n", s.refreshRate)
			s.refreshRate = 5 * time.Second
		}

		s.p.PrintColored(printer.Blue, "Starting market analysis for %s (updates every %v)\n",
			strings.ToUpper(s.ticker), s.refreshRate)
		time.Sleep(1 * time.Second)
		s.p.Reset()
		s.p = s.p.CleanLine()
		watchList := monitor.NewWatchList(s.p, s.refreshRate)
		go watchList.StartMonitoring(func() error {
			return s.analyse(ctx, cli)
		})

		time.Sleep(s.lifespan)
		watchList.Stop()
	}

	return nil
}

func (s *analysisScope) analyse(ctx context.Context, cli api.Client) error {
	data, err := cli.GetOHLCV(ctx, s.ticker, &carnost.WithTimeframe{TimeFrame: carnost.TimeFrame(s.timeFrame)})
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("no data for ticker %s", s.ticker)
	}

	m := make(map[signals.Operation]int, len(s.cfg.Strategies))

	for _, strat := range s.cfg.Strategies {
		operation := strat.Strategy.Execute(data)
		m[operation]++
	}

	s.p.Reset()
	strategies.NewAnalysisInput(s.p.CleanLine(), s.cfg.Strategies...).GenerateAnalysis()
	s.p.Printf("Sentiment is:\n")
	s.printSentiment(signals.Buy, m)
	s.printSentiment(signals.Sell, m)
	s.printSentiment(signals.Setup, m)
	s.printSentiment(signals.NoOp, m)
	return nil
}

func (s *analysisScope) printSentiment(k signals.Operation, m map[signals.Operation]int) {
	var c printer.Color
	v, ok := m[k]
	if !ok {
		v = 0
	}
	switch k {
	case signals.Buy:
		c = printer.Green
	case signals.Sell:
		c = printer.Blue
	case signals.Setup:
		c = printer.Yellow
	case signals.NoOp:
		c = printer.White
	}
	key := printer.WrapInColor(strings.ToUpper(string(k)), c)
	s.p.Printf(key+":\t%2.f%%\n", float32(v)/float32(len(s.cfg.Strategies))*100)
}
