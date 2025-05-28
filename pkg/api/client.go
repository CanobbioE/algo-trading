package api

type Client interface {
	GetOHLCV(ticker string, opts ...CallOption) ([]*OHLCV, error)
	GetEOD(ticker string, opts ...CallOption) (*EOD, error)
}

type Options interface{}

type CallOption interface {
	Apply(Options)
}

type OHLCV struct {
	Timestamp       float64
	Open            float64
	High            float64
	Low             float64
	Close           float64
	WeightedAverage float64
	Volume          float64
}

type EOD struct {
	Ticker       string
	ISIN         string
	Opening      float64
	MaxToday     float64
	MinToday     float64
	CurrentPrice float64
}
