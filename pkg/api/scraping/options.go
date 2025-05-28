package scraping

import (
	"github.com/canobbioe/algo-trading/pkg/api"
	"github.com/canobbioe/algo-trading/pkg/utilities"
)

type callOptions struct {
	timeFrame  *string
	sampleTime *string
}

type TimeFrame string

const (
	Daily        TimeFrame = "1d"
	Monthly      TimeFrame = "1m"
	Quarterly    TimeFrame = "3m"
	HalfYearly   TimeFrame = "6m"
	Yearly       TimeFrame = "1y"
	Triennial    TimeFrame = "3y"
	Quinquennial TimeFrame = "5y"
)

type WithTimeframe struct {
	TimeFrame TimeFrame
}

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
