package cmd

import (
	"encoding/json"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/printer"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

type monitorScope struct {
	p           printer.Printer
	cfg         *config.Config
	cfgFile     string
	refreshRate time.Duration
	lifespan    time.Duration
}

func (s *monitorScope) preRunE(_ *cobra.Command, _ []string) error {
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

func (s *monitorScope) runE(cmd *cobra.Command, _ []string) error {
	scanner := monitor.NewMarketScanner(s.cfg.Strategies, s.cfg.StockUniverse, s.cfg.Filters, scraping.NewClient(), s.p)

	watchList := monitor.NewWatchList(s.p, s.refreshRate)
	s.p.PrintColored(printer.Blue, "Starting market monitoring (updates every %v)\n", s.refreshRate)
	go watchList.StartMonitoring(func() error {
		scores, err := scanner.ScanMarket(cmd.Context())
		if err != nil {
			return err
		}

		if len(scores) == 0 {
			s.p.PrintColored(printer.Red, "No stocks meet current criteria\n")
			return nil
		}
		s.p.Printf("Found %d opportunities\n", len(scores))
		// show top 5 stock picks (todo: make it configurable)
		scanner.GenerateReport(scores, 5)
		s.checkAlerts(scores)
		return nil
	})

	time.Sleep(s.lifespan)
	watchList.Stop()
	return nil
}

// checkAlerts sends alerts for exceptional opportunities.
func (s *monitorScope) checkAlerts(scores []*monitor.StockScore) {
	// TODO: send me an email
	for _, score := range scores {
		if score.Opportunity >= monitor.OpportunityMedium && score.Risk <= monitor.RiskMedium {
			s.p.PrintColored(printer.Green, "HIGH OPPORTUNITY ALERT: %s (Score: %.2f, Confidence: %.1f%%)\n",
				score.Symbol, score.WeightedScore, score.Confidence*100)
		}
	}
}

func init() {
	s := &monitorScope{
		p: &printer.Standard{},
	}
	monitorCmd := &cobra.Command{
		Use:     "monitor",
		Short:   "Monitor the market",
		Long:    "Monitor the market.",
		PreRunE: s.preRunE,
		RunE:    s.runE,
	}

	monitorCmd.Flags().StringVarP(&s.cfgFile, "config", "c", "", "Path to config file")
	monitorCmd.Flags().DurationVarP(&s.refreshRate, "refresh", "r", 10*time.Minute, "Scan refresh rate")
	monitorCmd.Flags().DurationVarP(&s.lifespan, "life", "l", 1*time.Hour, "How long the monitor should run for")

	utilities.Must(monitorCmd.MarkFlagRequired("config"))
	rootCmd.AddCommand(monitorCmd)
}
