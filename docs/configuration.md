# Market Scanner Configuration Documentation

This document explains each configuration parameter used in the market scanner system and how they affect the analysis and filtering of investment opportunities.

## Table of Contents
- [Stock Universe](#stock-universe)
- [Trading Strategies](#trading-strategies)
- [Strategy Thresholds](#strategy-thresholds)
- [General Parameters](#general-parameters)
- [Scan Filters](#scan-filters)
- [Configuration Examples](#configuration-examples)

---

## Stock Universe

**Purpose**: Defines the list of financial instruments to analyze during market scans.

```json
"stock_universe": [
  "ENI.MTA",
  "A2A.MTA", 
  "UNI.MTA",
  "LDO.MTA",
  "QQQS.ETF",
  "QQQ3.ETF",
  "1GOOGL.MTA",
  "1MSFT.MTA"
]
```

### Parameters:
- **Format**: Array of ticker symbols with exchange suffixes
- **Exchange Suffixes**:
    - `.MTA` = Milan Stock Exchange (Mercato Telematico Azionario)
    - `.ETF` = Exchange-Traded Fund
    - Prefix `1` = Fractional shares or specific listing variants

### Usage Guidelines:
- **Diversification**: Include stocks from different sectors and asset classes
- **Liquidity**: Ensure included instruments have sufficient trading volume
- **Geographic**: Consider different markets and currencies for global exposure

---

## Trading Strategies

**Purpose**: Defines which trading strategies to use and their relative importance in decision-making.

```json
"strategies": [
  {
    "strategy": "BREAKOUT",
    "weight": 1.0
  },
  {
    "strategy": "VWAP", 
    "weight": 1.0
  },
  {
    "strategy": "MEANREVERSION",
    "weight": 1.0
  },
  {
    "strategy": "BOLLINGER",
    "weight": 1.0
  }
]
```

### Available Strategies:

#### BREAKOUT
- **Purpose**: Identifies stocks breaking above resistance levels
- **Best For**: Trending markets, momentum plays
- **Signals**: BUY when price breaks above recent highs with volume confirmation

#### VWAP (Volume Weighted Average Price)
- **Purpose**: Compares current price to volume-weighted average
- **Best For**: Intraday trading, institutional benchmarking
- **Signals**: BUY when price is below VWAP (potential value), SELL when significantly above

#### MEANREVERSION
- **Purpose**: Identifies stocks that have moved too far from their average price
- **Best For**: Range-bound markets, oversold/overbought conditions
- **Signals**: BUY when price is significantly below moving average

#### BOLLINGER
- **Purpose**: Uses Bollinger Bands to identify volatility and price extremes
- **Best For**: Volatility-based trading, squeeze patterns
- **Signals**: BUY on band squeezes or bounces off lower band

### Weight Configuration:
- **Range**: 0.1 to 5.0 (typically)
- **Equal Weights (1.0)**: All strategies have equal influence
- **Higher Weights (>1.0)**: Strategy has more influence on final decision
- **Lower Weights (<1.0)**: Strategy has less influence on final decision

**Example Weight Scenarios**:
```json
// Conservative approach - favor mean reversion
{"strategy": "MEANREVERSION", "weight": 1.5}
{"strategy": "BREAKOUT", "weight": 0.8}

// Aggressive approach - favor breakouts  
{"strategy": "BREAKOUT", "weight": 1.8}
{"strategy": "MEANREVERSION", "weight": 0.6}
```

---

## Strategy Thresholds

**Purpose**: Fine-tune the sensitivity and behavior of individual trading strategies.

```json
"thresholds": {
  "atr_period": 3,
  "low_atr_threshold": 0.05, 
  "high_atr_threshold": 0.2,
  "low_lookback": 10,
  "high_lookback": 3,
  "volume_threshold": 1.2,
  "deviation": 0.02,
  "squeeze": 0.1
}
```

### ATR (Average True Range) Parameters:

#### `atr_period`
- **Purpose**: Number of periods to calculate Average True Range
- **Range**: 2-20 periods
- **Lower Values**: More sensitive to recent volatility changes
- **Higher Values**: Smoother, less reactive to short-term spikes
- **Default**: 3 (responsive to recent volatility)

#### `low_atr_threshold`
- **Purpose**: Minimum volatility threshold (as percentage)
- **Range**: 0.01-0.10 (1%-10%)
- **Usage**: Filters out stocks with very low volatility
- **Impact**: Higher values = only high-volatility stocks considered
- **Default**: 0.05 (5% minimum daily volatility)

#### `high_atr_threshold`
- **Purpose**: Maximum volatility threshold (as percentage)
- **Range**: 0.10-1.00 (10%-100%)
- **Usage**: Filters out extremely volatile stocks
- **Impact**: Lower values = excludes high-risk volatile stocks
- **Default**: 0.2 (20% maximum daily volatility)

### Lookback Parameters:

#### `low_lookback`
- **Purpose**: Long-term trend analysis period
- **Range**: 5-50 periods
- **Usage**: Determines longer-term moving averages and trends
- **Impact**: Higher values = more stable, longer-term signals
- **Default**: 10 (two-week trend for daily data)

#### `high_lookback`
- **Purpose**: Short-term trend analysis period
- **Range**: 2-10 periods
- **Usage**: Determines short-term moving averages and signals
- **Impact**: Lower values = more responsive to recent price action
- **Default**: 3 (three-day trend for daily data)

### Volume and Signal Parameters:

#### `volume_threshold`
- **Purpose**: Volume multiplier for breakout confirmation
- **Range**: 1.0-3.0
- **Usage**: Volume must be X times average volume for valid breakout
- **Impact**: Higher values = requires stronger volume confirmation
- **Default**: 1.2 (20% above average volume)

#### `deviation`
- **Purpose**: Mean reversion sensitivity (as percentage)
- **Range**: 0.01-0.10 (1%-10%)
- **Usage**: How far price must deviate from mean to trigger signal
- **Impact**: Lower values = more sensitive mean reversion signals
- **Default**: 0.02 (2% deviation from mean)

#### `squeeze`
- **Purpose**: Bollinger Band squeeze threshold
- **Range**: 0.05-0.30
- **Usage**: Minimum band width ratio to identify squeezes
- **Impact**: Lower values = more sensitive squeeze detection
- **Default**: 0.1 (bands must be within 10% of normal width)

---

## General Parameters

### `lookback`
```json
"lookback": 3
```
- **Purpose**: Default lookback period for strategies that don't specify their own
- **Range**: 2-20 periods
- **Usage**: Moving averages, trend calculations, volatility measures
- **Impact**: Higher values = smoother but less responsive signals

### `bollinger_coefficient`
```json
"bollinger_coefficient": 2.0
```
- **Purpose**: Standard deviation multiplier for Bollinger Bands
- **Range**: 1.0-3.0
- **Standard Values**:
    - `1.0`: Very tight bands (68% of price action)
    - `2.0`: Standard bands (95% of price action)
    - `2.5`: Wide bands (99% of price action)
- **Impact**: Higher values = wider bands, fewer signals, higher confidence

---

## Scan Filters

**Purpose**: Filter and rank stocks based on quality and risk criteria.

```json
"scan_filters": {
  "min_confidence": 0.5,
  "min_weightedScore": 0.5, 
  "max_risk": "HIGH",
  "min_opportunity": "LOW",
  "min_volume": 1000,
  "required_signals": 1
}
```

### Signal Quality Filters:

#### `min_confidence`
- **Purpose**: Minimum percentage of strategies that must agree
- **Range**: 0.0-1.0 (0%-100%)
- **Usage**: 0.5 = at least 50% of strategies must signal BUY
- **Impact**: Higher values = fewer but higher-quality signals

#### `min_weightedScore`
- **Purpose**: Minimum weighted score across all strategies
- **Range**: 0.0-5.0 (depends on strategy weights)
- **Usage**: Accounts for strategy weights in decision making
- **Impact**: Higher values = only strongest signals pass filter

#### `required_signals`
- **Purpose**: Minimum number of BUY signals required
- **Range**: 1 to number of strategies
- **Usage**: Absolute minimum regardless of weights/confidence
- **Impact**: Higher values = require broader strategy agreement

### Risk and Opportunity Filters:

#### `max_risk`
- **Purpose**: Maximum acceptable risk level
- **Options**: `"LOW"`, `"MEDIUM"`, `"HIGH"`
- **Usage**: Filters out stocks above specified risk tolerance
- **Risk Factors**: Volatility, volume, sector

#### `min_opportunity`
- **Purpose**: Minimum opportunity level required
- **Options**: `"LOW"`, `"MEDIUM"`, `"HIGH"`
- **Usage**: Only includes stocks with sufficient upside potential
- **Opportunity Factors**: Signal strength, momentum, technical setup

### Market Filters:

#### `min_volume`
- **Purpose**: Minimum average daily trading volume
- **Range**: 100-1,000,000+ shares
- **Usage**: Ensures adequate liquidity for trading
- **Impact**: Higher values = exclude thinly traded stocks

---

## Configuration Examples
Three sample configurations are provided in the sample-configs folder:

### [Conservative Growth Configuration](../sample-configs/conservative.json)
> Focuses on stable, lower-risk opportunities with a preference for range-bound or mean-reverting behavior and well-capitalized companies.

### [Aggressive Growth Configuration](../sample-configs/aggressive.json)
> Focuses on breakout and momentum strategies in high-volatility environments with looser filters for risk.

### [Value Investing Configuration](../sample-configs/value.json)
> Targets undervalued, stable stocks using VWAP and mean-reversion with high-confidence and low-risk filters.


---

## Tuning Guidelines

### For Bull Markets:
- Increase BREAKOUT strategy weight
- Lower `min_confidence` requirements
- Accept higher risk levels
- Increase `volume_threshold` for breakouts

### For Bear Markets:
- Increase MEANREVERSION strategy weight
- Higher `min_confidence` requirements
- Lower maximum risk tolerance

### For Volatile Markets:
- Increase BOLLINGER strategy weight
- Adjust `atr_threshold` ranges appropriately
- Lower `bollinger_coefficient` for tighter bands
- Require higher volume confirmation

### For Range-Bound Markets:
- Emphasize MEANREVERSION and VWAP strategies
- Reduce BREAKOUT strategy weight
- Tighter `deviation` thresholds
- Higher confidence requirements