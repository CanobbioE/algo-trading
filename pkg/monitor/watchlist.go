package monitor

import (
	"fmt"
	"time"

	"github.com/CanobbioE/algo-trading/pkg/printer"
)

// WatchList maintains a dynamic watch list based on scan results.
type WatchList struct {
	p              printer.Printer
	stopChan       chan bool
	updateInterval time.Duration
}

// NewWatchList creates a new watch list.
func NewWatchList(p printer.Printer, updateInterval time.Duration) *WatchList {
	return &WatchList{
		updateInterval: updateInterval,
		stopChan:       make(chan bool),
		p:              p,
	}
}

// StartMonitoring begins continuous market monitoring.
func (wl *WatchList) StartMonitoring(callback func() error) {
	ticker := time.NewTicker(wl.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			wl.p.PrintColored(printer.Blue, "Running market scan...\n")
			if callback != nil {
				// todo: maybe return a bool "shouldContinue" to improve err handling.
				if err := callback(); err != nil {
					wl.p.PrintColored(printer.Red, fmt.Sprintf("Scan error: %v\n", err))
					continue
				}
			}
		case <-wl.stopChan:
			wl.p.PrintColored(printer.Yellow, "Stopping market monitoring")
			return
		}
	}
}

// Stop stops the monitoring.
func (wl *WatchList) Stop() {
	wl.stopChan <- true
}
