# Available commands

## analyse

Prints an analysis for a given ticker and the overall sentiment:

```shell
Key Technical Signals:
======================
- Latest Close (620.000) is equal to VWAP (620.000)
        -> No intraday edge for bulls or bears
- Latest Close (620.000) is equal to SMA (620.000)
        -> Trend indecision, can act as a magnet
- Deviation from SMA: 0.00%
        -> trend strength: near neutral
- Resistance Level: 620.000, Support Level: 620.000
- Price is currently inside breakout range
- Volume is below average
        -> Watch for confirmation before acting- Relative Strength: 50.00
        -> Neutral (sideway market)
- Bollinger Band SMA: 620.000, Upper: 620.000, Lower: 620.000, Width: 0.000
        -> Consider selling (price near or above upper band)
- ATR: 0.0000
        -> Volatility is low
======================
Sentiment is:
BUY:     0%
SELL:    0%
SETUP:  25%
NOOP:   75%
```

**Supported Flags:**

| Shorthand | Full Name   | Type       | Description                                                           | Default   |
|-----------|-------------|------------|-----------------------------------------------------------------------|-----------|
| -c        | --config    | [string]   | path to config file (required)                                        |           |
| -t        | --ticker    | [string]   | Stock ticker to use (required)                                        |           |
| -f        | --timeframe | [string]   | Time frame to use (one of [`1d`, `1m`, `3m`, `6m`, `1y`, `3y`, `5y`]) | `1d`      |
| -l        | --life      | [duration] | How long the monitor should run for                                   | `1h0m0s`  |
| -r        | --refresh   | [duration] | Scan refresh rate (min `5s`)                                          | `10m0s`   |
| -m        | --mode      | [string]   | How the command will run (one of `continue` or `onetime`)             | `onetime` |

## monitor

Continuously monitor the market for profitable trades

```shell
Starting market monitoring (updates every 10s)
Running market scan...
Scanning 589 stocks...
Filtering 589 results...
Found 215 opportunities

=== MARKET SCAN RESULTS ===
Found 215 stocks meeting criteria
Showing top 5 opportunities:

Rank #1: RAD.MTA
----------------
  Price: $1.08
  Signals: 1 BUY, 0 SELL, 2 HOLD, 1 SETUP
  Confidence: 37.5%
  Weighted Score: 1.00
  Risk: High | Opportunity: Low
  Volume: 1000
  Reasoning: *strategies.VWAPStrategy suggests BUY
...
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
Scanning 589 stocks...
Filtering 589 results...

=== MARKET SCAN RESULTS ===
Found 215 stocks meeting criteria
Showing top 10 opportunities:

Rank #1: INC.MTA
----------------
  Price: $0.18
  Signals: 1 BUY, 0 SELL, 2 HOLD, 1 SETUP
  Confidence: 37.5%
  Weighted Score: 1.00
  Risk: High | Opportunity: Low
  Volume: 128
  Reasoning: *strategies.VWAPStrategy suggests BUY
...
```

**Supported Flags:**

| Shorthand | Full Name | Type     | Description                   | Default |
|-----------|-----------|----------|-------------------------------|---------|
| -c        | --config  | [string] | path to config file (required |         |
