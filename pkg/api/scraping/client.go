package scraping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

const (
	origin            = "https://charts.borsaitaliana.it"
	getInfoEndpoint   = "/charts/services/ChartWService.asmx/GetInfos"
	getPricesEndpoint = "/charts/services/ChartWService.asmx/GetPricesWithVolume"
)

type scrapingClient struct {
	client *http.Client
}

// NewClient creates a new scraping client that implements api.Client.
func NewClient() api.Client {
	return &scrapingClient{
		client: http.DefaultClient,
	}
}

// GetEOD returns End Of Day data.
func (c *scrapingClient) GetEOD(ctx context.Context, ticker string, opts ...api.CallOption) (*api.EOD, error) {
	o := &callOptions{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	var data bytes.Buffer
	err := json.NewEncoder(&data).Encode(newGetEODRequest(ticker, o))
	if err != nil {
		return nil, fmt.Errorf("failed to encode EOD request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, origin+getInfoEndpoint, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to create EOD request: %w", err)
	}
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", origin)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send EOD request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	d := json.NewDecoder(resp.Body)
	var out eodResponse
	err = d.Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("failed to decode EOD response: %w", err)
	}

	/*
		[
		  "ENI", // ticker
		  "IT0003132476", // isin
		  12.87, // opening
		  12.904, // max today
		  12.574, // min today
		  12.756 // current price
		]
	*/
	return &api.EOD{
		Ticker:       out.D[0].(string),
		ISIN:         out.D[1].(string),
		Opening:      out.D[2].(float64),
		MaxToday:     out.D[3].(float64),
		MinToday:     out.D[4].(float64),
		CurrentPrice: out.D[5].(float64),
	}, nil
}

// GetOHLCV returns Open, High, Low, Close, Volume data.
func (c *scrapingClient) GetOHLCV(ctx context.Context, ticker string, opts ...api.CallOption) ([]*api.OHLCV, error) {
	o := &callOptions{}
	for _, opt := range opts {
		opt.Apply(o)
	}

	var data bytes.Buffer
	err := json.NewEncoder(&data).Encode(newGetOCHLVRequest(ticker, o))
	if err != nil {
		return nil, fmt.Errorf("failed to encode OHLCV request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, origin+getPricesEndpoint, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to create OHLCV request: %w", err)
	}
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "en-US,en;q=0.9")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", origin)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send OCHLV request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	d := json.NewDecoder(resp.Body)
	var out getOCHLVResponse
	err = d.Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("failed to decode OCHLV response: %w", err)
	}

	/*
	   [
	       1716768000000, // timestamp
	       14.582, // open price
	       14.4, // high
	       14.582, // low
	       14.37, // close
	       14.582 // weighted average
	   ],
	*/
	response := make([]*api.OHLCV, len(out.D))
	for i, timeframe := range out.D {
		response[i] = &api.OHLCV{
			Timestamp:       utilities.SafeConvert[float64](timeframe[0]),
			Open:            utilities.SafeConvert[float64](timeframe[1]),
			High:            utilities.SafeConvert[float64](timeframe[2]),
			Low:             utilities.SafeConvert[float64](timeframe[3]),
			Close:           utilities.SafeConvert[float64](timeframe[4]),
			WeightedAverage: utilities.SafeConvert[float64](timeframe[5]),
			Volume:          utilities.SafeConvert[float64](timeframe[6]),
		}
	}

	return response, nil
}
