# algo-trading
My algo trading strategy


# Available commands
## analyse
Prints an analysis for a given ticker and the overall sentiment:
```shell
Dynamic Lookback Period: 6
Key Technical Signals:
======================
- Latest Close (128.900) is above VWAP (128.833)
        -> Suggests buyers are in control (BULLISH)
- Latest Close (128.900) is above SMA (128.800)
        -> Indicates uptrend continuation or strength (support in rising trend)
- Deviation from SMA: 0.00%
        -> trend strength: near neutral
- Resistance Level: 128.900, Support Level: 128.750
- Price is currently inside breakout range
- Volume is above average
        -> Validates potential breakout
- Relative Strength: 50.00
        -> Neutral (sideway market)
- Bollinger Band SMA: 128.800, Upper: 128.941, Lower: 128.659, Width: 0.283
        -> Consider selling (price near or above upper band)
- ATR: 0.0533
        -> Volatility is moderate or high
======================
Sentiment is:
NOOP:   75%
BUY:    25%

```

**Supported Flags:**

|Shorthand| Full Name| Description|
|---|---|---|
| -h | --help         |      help for analyse
|-t |--ticker    |   [string] Stock ticker to use (required)
|-f |--timeframe  |  [string]  Time frame to use (one of [`1d`, `1m`, `3m`, `6m`, `1y`, `3y`, `5y`] - default `1d`)

**Usage**
```
 go run . analyse -t VWCE.ETF -t 1y
```

  
