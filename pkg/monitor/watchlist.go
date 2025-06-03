package monitor

import (
	"fmt"
	"log"
	"time"
)

// WatchList maintains a dynamic watch list based on scan results
type WatchList struct {
	stocks         []*StockScore
	scanner        *MarketScanner
	updateInterval time.Duration
	stopChan       chan bool
}

// NewWatchList creates a new watch list
func NewWatchList(scanner *MarketScanner, updateInterval time.Duration) *WatchList {
	return &WatchList{
		scanner:        scanner,
		updateInterval: updateInterval,
		stopChan:       make(chan bool),
	}
}

// StartMonitoring begins continuous market monitoring
func (wl *WatchList) StartMonitoring() {
	ticker := time.NewTicker(wl.updateInterval)
	defer ticker.Stop()

	fmt.Printf("Starting market monitoring (updates every %v)\n", wl.updateInterval)

	for {
		select {
		case <-ticker.C:
			fmt.Printf("\n[%s] Running market scan...\n", time.Now().Format("15:04:05"))

			scores, err := wl.scanner.ScanMarket()
			if err != nil {
				log.Printf("Scan error: %v", err)
				continue
			}

			wl.stocks = scores

			if len(scores) > 0 {
				fmt.Printf("Found %d opportunities\n", len(scores))
				wl.scanner.GenerateReport(scores, 5) // Show top 5

				// Alert on high-opportunity stocks
				wl.checkAlerts(scores)
			} else {
				fmt.Println("No stocks meet current criteria")
			}

		case <-wl.stopChan:
			fmt.Println("Stopping market monitoring")
			return
		}
	}
}

// checkAlerts sends alerts for exceptional opportunities
func (wl *WatchList) checkAlerts(scores []*StockScore) {
	for _, score := range scores {
		if score.Opportunity == OpportunityHigh && score.Risk <= RiskMedium {
			fmt.Printf("🚨 HIGH OPPORTUNITY ALERT: %s (Score: %.2f, Confidence: %.1f%%)\n",
				score.Symbol, score.WeightedScore, score.Confidence*100)
		}
	}
}

// Stop stops the monitoring
func (wl *WatchList) Stop() {
	wl.stopChan <- true
}
