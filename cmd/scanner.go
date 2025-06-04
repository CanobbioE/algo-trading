package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
	"github.com/spf13/cobra"
	"os"
)

type scanScope struct {
	cfgFile string
	cfg     *config.Config
}

func (s *scanScope) preRunE(_ *cobra.Command, _ []string) error {
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

func (s *scanScope) runE(_ *cobra.Command, _ []string) error {
	scanner := monitor.NewMarketScanner(s.cfg.Strategies, s.cfg.StockUniverse, s.cfg.Filters, scraping.NewClient())

	fmt.Println("=== ONE-TIME MARKET SCAN ===")
	scores, err := scanner.ScanMarket()
	if err != nil {
		return err
	}

	scanner.GenerateReport(scores, 10)
	return nil
}

func init() {
	s := &scanScope{}
	scanCmd := &cobra.Command{
		Use:     "scan",
		Short:   "Scan the market once",
		Long:    "Scan the market once.",
		PreRunE: s.preRunE,
		RunE:    s.runE,
	}

	scanCmd.Flags().StringVarP(&s.cfgFile, "config", "c", "", "path to config file")

	utilities.Must(scanCmd.MarkFlagRequired("config"))
	rootCmd.AddCommand(scanCmd)
}
