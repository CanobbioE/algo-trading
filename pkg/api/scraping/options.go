package scraping

import (
	"github.com/CanobbioE/algo-trading/pkg/api"
	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

type callOptions struct {
	timeFrame  *string
	sampleTime *string
}

// TimeFrame represents a predefined period of time.
type TimeFrame string

const (
	// Daily is a TimeFrame representing one day.
	Daily TimeFrame = "1d"
	// Monthly is a TimeFrame representing one month.
	Monthly TimeFrame = "1m"
	// Quarterly is a TimeFrame representing three months.
	Quarterly TimeFrame = "3m"
	// HalfYearly is a TimeFrame representing six months.
	HalfYearly TimeFrame = "6m"
	// Yearly is a TimeFrame representing one year.
	Yearly TimeFrame = "1y"
	// Triennial is a TimeFrame representing  three years.
	Triennial TimeFrame = "3y"
	// Quinquennial is a TimeFrame representing five years.
	Quinquennial TimeFrame = "5y"
)

// WithTimeframe applies a timeframe option to the API requests.
type WithTimeframe struct {
	TimeFrame TimeFrame
}

// Apply the WithTimeframe option, it automatically sets the sample time based on the value of the timeframe.
func (o *WithTimeframe) Apply(in api.Options) {
	opts, ok := in.(*callOptions)
	if !ok {
		return
	}
	opts.timeFrame = utilities.ToPointer(string(o.TimeFrame))
	switch o.TimeFrame {
	case Daily:
		opts.sampleTime = utilities.ToPointer("1mm")
	case Monthly, Quarterly, HalfYearly, Yearly, Triennial, Quinquennial:
		opts.sampleTime = utilities.ToPointer("1d")
	}
}
