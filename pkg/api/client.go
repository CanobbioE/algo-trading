package api

import "context"

// Client is the generic API client interface all concrete API clients must implement.
type Client interface {
	// GetOHLCV returns Open, High, Low, Close, Volume data.
	GetOHLCV(ctx context.Context, ticker string, opts ...CallOption) ([]*OHLCV, error)
	// GetEOD returns End Of Day data.
	GetEOD(ctx context.Context, ticker string, opts ...CallOption) (*EOD, error)
}

// Options is any collection of options.
type Options any

// CallOption is the interface that any API request option must implement.
type CallOption interface {
	// Apply the current CallOption to the options.
	Apply(o Options)
}

// OHLCV represents the standard response for Open, High, Low, Close, Volume data.
type OHLCV struct {
	Timestamp       float64
	Open            float64
	High            float64
	Low             float64
	Close           float64
	WeightedAverage float64
	Volume          float64
}

// EOD represents the standard response for End Of Day data.
type EOD struct {
	Ticker       string
	ISIN         string
	Opening      float64
	MaxToday     float64
	MinToday     float64
	CurrentPrice float64
}
