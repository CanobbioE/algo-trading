package cmd

import (
	"encoding/json"
	"os"

	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/config"
	"github.com/CanobbioE/algo-trading/pkg/monitor"
	"github.com/CanobbioE/algo-trading/pkg/printer"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
	"github.com/spf13/cobra"
)

type scanScope struct {
	p       printer.Printer
	cfg     *config.Config
	cfgFile string
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

func (s *scanScope) runE(cmd *cobra.Command, _ []string) error {
	scanner := monitor.NewMarketScanner(s.cfg.Strategies, s.cfg.StockUniverse, s.cfg.Filters, scraping.NewClient(), s.p)

	s.p.Printf("=== ONE-TIME MARKET SCAN ===\n")
	scores, err := scanner.ScanMarket(cmd.Context())
	if err != nil {
		return err
	}

	scanner.GenerateReport(scores, 10)
	return nil
}

func init() {
	s := &scanScope{
		p: &printer.Standard{},
	}
	scanCmd := &cobra.Command{
		Use:     "scan",
		Short:   "Scan the market once",
		Long:    "Scan the market once.",
		PreRunE: s.preRunE,
		RunE:    s.runE,
	}

	scanCmd.Flags().StringVarP(&s.cfgFile, "config", "c", "", "Path to config file")

	utilities.Must(scanCmd.MarkFlagRequired("config"))
	rootCmd.AddCommand(scanCmd)
}
