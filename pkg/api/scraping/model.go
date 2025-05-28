package scraping

import "github.com/canobbioe/algo-trading/pkg/utilities"

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
	SampleTime           string `json:"SampleTime"`
	TimeFrame            string `json:"TimeFrame"`
	RequestedDataSetType string `json:"RequestedDataSetType"`
	ChartPriceType       string `json:"ChartPriceType"`
	Key                  string `json:"Key"`
	OffSet               int    `json:"OffSet"`
	FromDate             any    `json:"FromDate"`
	ToDate               any    `json:"ToDate"`
	UseDelay             bool   `json:"UseDelay"`
	KeyType              string `json:"KeyType"`
	KeyType2             string `json:"KeyType2"`
	Language             string `json:"Language"`
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
			Key:                  key,
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

// TODO
func newGetEODRequest(key string) *getEODReq {
	return &getEODReq{
		Request: &getEODInnerReq{
			Key:                  "",
			TimeFrame:            "",
			SampleTime:           "",
			RequestedDataSetType: "",
			ChartPriceType:       "",
			UseDelay:             "",
			KeyType:              "",
		},
	}
}
