package monitor

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/api/scraping"
	"github.com/CanobbioE/algo-trading/pkg/signals"
	"github.com/CanobbioE/algo-trading/pkg/strategies"
)

// StockScore represents the analysis result for a single stock
type StockScore struct {
	Symbol        string
	BuySignals    int
	SellSignals   int
	HoldSignals   int
	SetupSignals  int
	WeightedScore float64
	Confidence    float64
	Reasoning     []string
	LastPrice     float64
	Volume        float64
	MarketCap     float64
	Risk          RiskLevel
	Opportunity   OpportunityLevel
}

// MarketScanner scans multiple stocks and ranks them
type MarketScanner struct {
	strategies     []*strategies.StrategyWeight
	client         api.Client
	stockUniverse  []string
	filters        *ScanFilters
	maxConcurrency int
}

// NewMarketScanner creates a new market scanner
func NewMarketScanner(strategies []*strategies.StrategyWeight, stockList []string, filters *ScanFilters, cli api.Client) *MarketScanner {
	return &MarketScanner{
		strategies:     strategies,
		client:         cli,
		stockUniverse:  stockList,
		maxConcurrency: 10, // Limit concurrent API calls
		filters:        filters,
	}
}

// ScanMarket analyzes all stocks in the universe
func (ms *MarketScanner) ScanMarket() ([]*StockScore, error) {
	fmt.Printf("Scanning %d stocks...\n", len(ms.stockUniverse))

	// Channel to control concurrency
	semaphore := make(chan struct{}, ms.maxConcurrency)
	results := make(chan *StockScore, len(ms.stockUniverse))
	errors := make(chan error, len(ms.stockUniverse))

	var wg sync.WaitGroup

	// Analyze each stock concurrently
	for _, symbol := range ms.stockUniverse {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			score, err := ms.analyzeStock(sym)
			if err != nil {
				errors <- fmt.Errorf("error analyzing %s: %v", sym, err)
				return
			}

			results <- score
		}(symbol)
	}

	// Close channels when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// Collect results
	var scores []*StockScore
	var scanErrors []error

	for {
		select {
		case score, ok := <-results:
			if !ok {
				results = nil
			} else {
				scores = append(scores, score)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				scanErrors = append(scanErrors, err)
			}
		}

		if results == nil && errors == nil {
			break
		}
	}

	// Log errors but don't fail the entire scan
	for _, err := range scanErrors {
		log.Printf("Scan error: %v", err)
	}

	// Filter and sort results
	filteredScores := ms.filterResults(scores)
	ms.sortByOpportunity(filteredScores)

	return filteredScores, nil
}

// analyzeStock performs strategy analysis on a single stock
func (ms *MarketScanner) analyzeStock(symbol string) (*StockScore, error) {
	// Get stock data
	data, err := ms.client.GetOHLCV(symbol, &scraping.WithTimeframe{TimeFrame: scraping.Daily})
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no data available for %s", symbol)
	}

	score := &StockScore{
		Symbol:    symbol,
		LastPrice: data[len(data)-1].Close,
		Volume:    data[len(data)-1].Volume,
	}

	signalCounts := make(map[signals.Operation]int)
	weightedScores := make(map[signals.Operation]float64)
	reasoning := make([]string, 0)

	for _, sw := range ms.strategies {
		operation := sw.Strategy.Execute(data)
		signalCounts[operation]++
		weightedScores[operation] += sw.Weight

		if operation == signals.Buy {
			reasoning = append(reasoning, fmt.Sprintf("%T suggests BUY", sw.Strategy))
		}
	}

	score.BuySignals = signalCounts[signals.Buy]
	score.SellSignals = signalCounts[signals.Sell]
	score.HoldSignals = signalCounts[signals.NoOp]
	score.SetupSignals = signalCounts[signals.Setup]
	score.WeightedScore = weightedScores[signals.Buy] - weightedScores[signals.Sell]
	score.Confidence = float64(score.BuySignals) / float64(len(ms.strategies))
	score.Reasoning = reasoning

	score.Risk = ms.calculateRisk(data, score)
	score.Opportunity = ms.calculateOpportunity(score)
	score.MarketCap = ms.estimateMarketCap(score.LastPrice, score.Volume)

	return score, nil
}

// calculateRisk assesses the risk level of a stock
func (ms *MarketScanner) calculateRisk(data []*api.OHLCV, score *StockScore) RiskLevel {
	if len(data) < 20 {
		return RiskHigh
	}

	// Calculate volatility (simplified)
	var priceChanges []float64
	for i := 1; i < len(data) && i < 21; i++ {
		change := (data[i].Close - data[i-1].Close) / data[i-1].Close
		priceChanges = append(priceChanges, change*change)
	}

	var avgVolatility float64
	for _, change := range priceChanges {
		avgVolatility += change
	}
	avgVolatility = avgVolatility / float64(len(priceChanges))
	volatility := avgVolatility * 100 // Convert to percentage

	// Risk assessment based on volatility and other factors
	switch {
	case volatility > 5.0 || score.Volume < 50000:
		return RiskHigh
	case volatility > 2.0 || score.Volume < 200000:
		return RiskMedium
	default:
		return RiskLow
	}
}

// calculateOpportunity assesses the opportunity level
func (ms *MarketScanner) calculateOpportunity(score *StockScore) OpportunityLevel {
	// Opportunity based on signal strength and confidence
	opportunityScore := score.WeightedScore + (score.Confidence * 2)

	switch {
	case opportunityScore >= 3.0 && score.BuySignals >= 3:
		return OpportunityHigh
	case opportunityScore >= 1.5 && score.BuySignals >= 2:
		return OpportunityMedium
	default:
		return OpportunityLow
	}
}

// estimateMarketCap provides a rough market cap estimate
func (ms *MarketScanner) estimateMarketCap(price, volume float64) float64 {
	// This is a very rough estimate - in practice, you'd get this from your data provider
	// Assuming average shares outstanding based on volume patterns
	estimatedShares := volume * 50 // Very rough heuristic
	return price * estimatedShares
}

// filterResults applies filters to the scan results
func (ms *MarketScanner) filterResults(scores []*StockScore) []*StockScore {
	var filtered []*StockScore

	for _, score := range scores {
		if score.Confidence >= ms.filters.MinConfidence &&
			score.WeightedScore >= ms.filters.MinWeightedScore &&
			score.Risk <= ms.filters.MaxRisk &&
			score.Opportunity >= ms.filters.MinOpportunity &&
			score.Volume >= ms.filters.MinVolume &&
			score.MarketCap >= ms.filters.MinMarketCap &&
			score.MarketCap <= ms.filters.MaxMarketCap &&
			score.BuySignals >= ms.filters.RequiredSignals {
			filtered = append(filtered, score)
		}
	}

	return filtered
}

// sortByOpportunity sorts stocks by opportunity score (best first)
func (ms *MarketScanner) sortByOpportunity(scores []*StockScore) {
	sort.Slice(scores, func(i, j int) bool {
		// Primary sort: Opportunity level
		if scores[i].Opportunity != scores[j].Opportunity {
			return scores[i].Opportunity > scores[j].Opportunity
		}

		// Secondary sort: Weighted score
		if scores[i].WeightedScore != scores[j].WeightedScore {
			return scores[i].WeightedScore > scores[j].WeightedScore
		}

		// Tertiary sort: Confidence
		return scores[i].Confidence > scores[j].Confidence
	})
}

// GenerateReport creates a formatted report of the top opportunities
func (ms *MarketScanner) GenerateReport(scores []*StockScore, topN int) {
	fmt.Println("\n=== MARKET SCAN RESULTS ===")
	fmt.Printf("Found %d stocks meeting criteria\n", len(scores))
	fmt.Printf("Showing top %d opportunities:\n\n", min(topN, len(scores)))

	for i, score := range scores {
		if i >= topN {
			break
		}

		fmt.Printf("Rank #%d: %s\n", i+1, score.Symbol)
		fmt.Printf("  Price: $%.2f\n", score.LastPrice)
		fmt.Printf("  Signals: %d BUY, %d SELL, %d HOLD\n",
			score.BuySignals, score.SellSignals, score.HoldSignals)
		fmt.Printf("  Confidence: %.1f%%\n", score.Confidence*100)
		fmt.Printf("  Weighted Score: %.2f\n", score.WeightedScore)
		fmt.Printf("  Risk: %v | Opportunity: %v\n", score.Risk, score.Opportunity)
		fmt.Printf("  Volume: %.0f\n", score.Volume)
		fmt.Printf("  Est. Market Cap: $%.2fM\n", score.MarketCap/1000000)

		if len(score.Reasoning) > 0 {
			fmt.Printf("  Reasoning: %s\n", score.Reasoning[0])
		}
		fmt.Println()
	}
}
