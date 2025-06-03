package cmd

import (
	"encoding/json"
	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
	"github.com/spf13/cobra"
	"os"
	"time"
)

type monitorScope struct {
	cfgFile     string
	cfg         *config.Config
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

func (s *monitorScope) runE(_ *cobra.Command, _ []string) error {
	scanner := monitor.NewMarketScanner(s.cfg.Strategies, s.cfg.StockUniverse, s.cfg.Filters, scraping.NewClient())

	watchList := monitor.NewWatchList(scanner, s.refreshRate)
	go watchList.StartMonitoring()

	time.Sleep(s.lifespan)
	watchList.Stop()
	return nil
}

func init() {
	s := &monitorScope{}
	monitorCmd := &cobra.Command{
		Use:     "monitor",
		Short:   "Monitor the market",
		Long:    "Monitor the market.",
		PreRunE: s.preRunE,
		RunE:    s.runE,
	}

	monitorCmd.Flags().StringVarP(&s.cfgFile, "config", "c", "", "path to config file")
	monitorCmd.Flags().DurationVarP(&s.refreshRate, "refresh", "r", 10*time.Minute, "scan refresh rate")
	monitorCmd.Flags().DurationVarP(&s.lifespan, "life", "l", 1*time.Hour, "how long the monitor should run for")

	utilities.Must(monitorCmd.MarkFlagRequired("config"))
	rootCmd.AddCommand(monitorCmd)
}
