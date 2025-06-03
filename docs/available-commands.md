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

| Shorthand | Full Name   | Type     | Description                                                           | Default |
|-----------|-------------|----------|-----------------------------------------------------------------------|---------|
| -t        | --ticker    | [string] | Stock ticker to use (required)                                        |         |
| -f        | --timeframe | [string] | Time frame to use (one of [`1d`, `1m`, `3m`, `6m`, `1y`, `3y`, `5y`]) | `1d`    |


## monitor

Continuously monitor the market for profitable trades

```shell
Starting market monitoring (updates every 10s)

Running market scan...
Scanning 8 stocks...
No stocks meet current criteria
```

**Supported Flags:**

| Shorthand | Full Name | Type       | Description                         | Default  |
|-----------|-----------|------------|-------------------------------------|----------|
| -c        | --config  | [string]   | path to config file (required)      |          |
| -l        | --life    | [duration] | how long the monitor should run for | `1h0m0s` |
| -r        | --refresh | [duration] | scan refresh rate                   | `10m0s`  |

## scan

Run a single scan of the market.

```shell
=== ONE-TIME MARKET SCAN ===
Scanning 8 stocks...


=== MARKET SCAN RESULTS ===
Found 0 stocks meeting criteria
Showing top 0 opportunities:
```

**Supported Flags:**

| Shorthand | Full Name | Type     | Description                   | Default |
|-----------|-----------|----------|-------------------------------|---------|
| -c        | --config  | [string] | path to config file (required |         |
