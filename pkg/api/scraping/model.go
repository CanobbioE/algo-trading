package scraping

import (
	"strings"

	"github.com/CanobbioE/algo-trading/pkg/utilities"
)

type eodResponse struct {
	D []any `json:"d"`
}

type getOCHLVResponse struct {
	D [][]any `json:"d"`
}

type getOCHLVReq struct {
	Request *getOCHLVInnerReq `json:"request"`
}

type getOCHLVInnerReq struct {
	FromDate             any    `json:"FromDate"`
	ToDate               any    `json:"ToDate"`
	SampleTime           string `json:"SampleTime"`
	TimeFrame            string `json:"TimeFrame"`
	RequestedDataSetType string `json:"RequestedDataSetType"`
	ChartPriceType       string `json:"ChartPriceType"`
	Key                  string `json:"Key"`
	KeyType              string `json:"KeyType"`
	KeyType2             string `json:"KeyType2"`
	Language             string `json:"Language"`
	OffSet               int    `json:"OffSet"`
	UseDelay             bool   `json:"UseDelay"`
}

type getEODReq struct {
	Request *getEODInnerReq `json:"request"`
}
type getEODInnerReq struct {
	Key                  string `json:"Key"`
	TimeFrame            string `json:"TimeFrame"`
	SampleTime           string `json:"SampleTime"`
	RequestedDataSetType string `json:"RequestedDataSetType"`
	ChartPriceType       string `json:"ChartPriceType"`
	UseDelay             string `json:"UseDelay"`
	KeyType              string `json:"KeyType"`
}

/*
	{
	    "SampleTime": "1mm",
	    "TimeFrame": "1d",
	    "RequestedDataSetType": "ohlc",
	    "ChartPriceType": "price",
	    "Key": "ENI.MTA",
	    "OffSet": 0,
	    "FromDate": null,
	    "ToDate": null,
	    "UseDelay": false,
	    "KeyType": "Topic",
	    "KeyType2": "Topic",
	    "Language": "it-IT"
	}
*/
func newGetOCHLVRequest(key string, options *callOptions) *getOCHLVReq {
	return &getOCHLVReq{
		Request: &getOCHLVInnerReq{
			Key:                  strings.ToUpper(key),
			SampleTime:           utilities.OptionalWithFallback(options.sampleTime, "1d"),
			TimeFrame:            utilities.OptionalWithFallback(options.timeFrame, "1mm"),
			RequestedDataSetType: "ohlc",
			ChartPriceType:       "price",
			OffSet:               0,
			FromDate:             nil, // todo: maybe we can play with this
			ToDate:               nil, // todo: maybe we can play with this
			UseDelay:             false,
			KeyType:              "Topic", // todo: not sure about this, test with more tickers
			KeyType2:             "Topic",
			Language:             "it-IT",
		},
	}
}

// TODO.
func newGetEODRequest(key string, options *callOptions) *getEODReq {
	return &getEODReq{
		Request: &getEODInnerReq{
			Key:                  key,
			SampleTime:           utilities.OptionalWithFallback(options.sampleTime, "1d"),
			TimeFrame:            utilities.OptionalWithFallback(options.timeFrame, "1mm"),
			RequestedDataSetType: "",
			ChartPriceType:       "",
			UseDelay:             "",
			KeyType:              "",
		},
	}
}
