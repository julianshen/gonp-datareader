# Data Sources Documentation

This document provides detailed information about each data source supported by gonp-datareader, including capabilities, limitations, symbol formats, API requirements, and rate limits.

---

## Table of Contents

1. [Yahoo Finance](#yahoo-finance)
2. [FRED](#fred-federal-reserve-economic-data)
3. [World Bank](#world-bank)
4. [Alpha Vantage](#alpha-vantage)
5. [Stooq](#stooq)
6. [IEX Cloud](#iex-cloud)
7. [Tiingo](#tiingo)
8. [OECD](#oecd)
9. [Eurostat](#eurostat)
10. [Comparison Matrix](#comparison-matrix)

---

## Yahoo Finance

### Overview
Yahoo Finance provides free access to historical stock market data including OHLCV (Open, High, Low, Close, Volume) data.

### API Key Required
**No** - Free access without registration

### Symbol Format
- **US Stocks**: `TICKER` (e.g., `AAPL`, `MSFT`, `GOOGL`)
- **International**: `TICKER.EXCHANGE` (e.g., `7203.T` for Toyota on Tokyo Exchange)
- **Indices**: `^INDEX` (e.g., `^GSPC` for S&P 500)
- **Crypto**: `SYMBOL-USD` (e.g., `BTC-USD`, `ETH-USD`)

### Data Available
- Daily OHLCV data
- Adjusted close prices
- Historical data (20+ years for most symbols)
- Real-time delayed quotes (15-20 minutes)

### Rate Limits
- No official rate limit
- Recommended: ~2000 requests/hour to avoid blocking
- Use with caching for production applications

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported

### Capabilities
- ✅ Free, no API key required
- ✅ Extensive historical data
- ✅ Global market coverage
- ✅ Cryptocurrency support
- ✅ Index data

### Limitations
- ❌ No real-time data (15-20 min delay)
- ❌ No fundamental data
- ❌ May block excessive requests
- ❌ No official API documentation

### Example Usage
```go
ctx := context.Background()
start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
end := time.Now()

// Single symbol
data, err := datareader.Read(ctx, "AAPL", "yahoo", start, end, nil)

// Multiple symbols (parallel)
symbols := []string{"AAPL", "MSFT", "GOOGL"}
dataMap, err := datareader.Read(ctx, symbols, "yahoo", start, end, nil)
```

### Links
- Website: https://finance.yahoo.com
- No official API documentation

---

## FRED (Federal Reserve Economic Data)

### Overview
FRED provides access to economic data from the Federal Reserve Bank of St. Louis, including thousands of economic indicators.

### API Key Required
**Optional** - Free API key available, some data accessible without

### Symbol Format
- **Series ID**: `SERIES_ID` (e.g., `GDP`, `UNRATE`, `DGS10`)
- Find series IDs on FRED website

### Data Available
- Economic indicators (GDP, unemployment, interest rates, etc.)
- Over 800,000 time series
- Historical data (varies by series, some >100 years)
- Updated regularly (daily to annually depending on series)

### Rate Limits
- **Without API key**: Limited requests
- **With API key**: 120 requests/minute
- Generous limits for research use

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported

### Capabilities
- ✅ Massive economic dataset (800K+ series)
- ✅ High-quality, authoritative data
- ✅ Free API key
- ✅ Extensive historical data
- ✅ Well-documented API

### Limitations
- ❌ Economic data only (no stock prices)
- ❌ API key recommended for production
- ❌ Rate limits apply

### Example Usage
```go
opts := &datareader.Options{
    APIKey: "your-fred-api-key", // Optional but recommended
}

// GDP data
data, err := datareader.Read(ctx, "GDP", "fred", start, end, opts)

// Multiple series
series := []string{"GDP", "UNRATE", "DGS10"}
dataMap, err := datareader.Read(ctx, series, "fred", start, end, opts)
```

### Links
- Website: https://fred.stlouisfed.org
- API Documentation: https://fred.stlouisfed.org/docs/api/fred/
- Get API Key: https://fred.stlouisfed.org/docs/api/api_key.html

---

## World Bank

### Overview
The World Bank provides international development indicators covering over 200 countries and regions.

### API Key Required
**No** - Free access without registration

### Symbol Format
- **Format**: `COUNTRY/INDICATOR`
- **Country**: ISO 3-letter code or special codes (e.g., `USA`, `CHN`, `all`, `WLD`)
- **Indicator**: World Bank indicator code (e.g., `NY.GDP.MKTP.CD`)
- **Example**: `USA/NY.GDP.MKTP.CD` (US GDP)

### Data Available
- Economic indicators (GDP, population, trade, etc.)
- Social indicators (education, health, poverty)
- Environmental indicators
- 1,400+ indicators across 200+ countries
- Historical data (1960-present for most indicators)

### Rate Limits
- No strict rate limits
- Paginated results (max 1000 per page)
- Recommended: Respect server with reasonable request rates

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported for multiple country/indicator combinations

### Capabilities
- ✅ Free, no API key required
- ✅ Comprehensive international data
- ✅ Authoritative source
- ✅ Well-structured API
- ✅ Long historical data

### Limitations
- ❌ Annual data only (no intra-year)
- ❌ Data lags (1-2 years for recent data)
- ❌ Limited to World Bank indicators
- ❌ Complex symbol format

### Example Usage
```go
// US GDP
data, err := datareader.Read(ctx, "USA/NY.GDP.MKTP.CD", "worldbank", start, end, nil)

// Multiple indicators
indicators := []string{
    "USA/NY.GDP.MKTP.CD",  // US GDP
    "CHN/NY.GDP.MKTP.CD",  // China GDP
    "WLD/SP.POP.TOTL",     // World population
}
dataMap, err := datareader.Read(ctx, indicators, "worldbank", start, end, nil)
```

### Links
- Website: https://data.worldbank.org
- API Documentation: https://datahelpdesk.worldbank.org/knowledgebase/articles/889392
- Indicator Catalog: https://data.worldbank.org/indicator

---

## Alpha Vantage

### Overview
Alpha Vantage provides real-time and historical stock market data, forex, and cryptocurrency data through a RESTful API.

### API Key Required
**Yes** - Free tier available (limited)

### Symbol Format
- **US Stocks**: `TICKER` (e.g., `AAPL`, `MSFT`)
- **International**: `TICKER.EXCHANGE` (e.g., `VOD.LON`)
- Case-insensitive

### Data Available
- Daily, weekly, monthly stock data
- Intraday data (1min, 5min, 15min, 30min, 60min)
- Adjusted OHLCV data
- Real-time quotes
- Historical data (20+ years)

### Rate Limits
- **Free tier**: 5 API calls/minute, 500 calls/day
- **Paid tiers**: Higher limits available
- Strict enforcement

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported (respects rate limits)

### Capabilities
- ✅ Real-time data
- ✅ Intraday data available
- ✅ High data quality
- ✅ Technical indicators
- ✅ Fundamental data available

### Limitations
- ❌ Strict rate limits (free tier)
- ❌ API key required
- ❌ Limited free tier usage
- ❌ Can be slow under free tier

### Example Usage
```go
opts := &datareader.Options{
    APIKey: "your-alphavantage-api-key",
}

data, err := datareader.Read(ctx, "AAPL", "alphavantage", start, end, opts)

// Multiple symbols (respects rate limits)
symbols := []string{"AAPL", "MSFT", "GOOGL"}
dataMap, err := datareader.Read(ctx, symbols, "alphavantage", start, end, opts)
```

### Links
- Website: https://www.alphavantage.co
- API Documentation: https://www.alphavantage.co/documentation
- Get API Key: https://www.alphavantage.co/support/#api-key

---

## Stooq

### Overview
Stooq provides free historical market data for stocks, indices, currencies, and commodities from global exchanges.

### API Key Required
**No** - Free access without registration

### Symbol Format
- **US Stocks**: `TICKER.US` (e.g., `AAPL.US`, `MSFT.US`)
- **Indices**: `^INDEX` (e.g., `^SPX`, `^DJI`)
- **International**: `TICKER.EXCHANGE` (e.g., `BMW.DE`, `7203.JP`)
- **Currencies**: `PAIR` (e.g., `EURUSD`, `GBPUSD`)

### Data Available
- Daily OHLCV data
- Historical data (10-20 years typically)
- Global market coverage
- Indices, stocks, forex, commodities

### Rate Limits
- No official rate limits
- Recommended: Reasonable request rates
- Use caching for production

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported

### Capabilities
- ✅ Free, no API key required
- ✅ International market coverage
- ✅ Simple CSV format
- ✅ Reliable data source
- ✅ Indices and forex included

### Limitations
- ❌ No real-time data (delayed)
- ❌ Limited to daily data
- ❌ No fundamental data
- ❌ Unofficial API (may change)

### Example Usage
```go
// Note the .US suffix for US stocks
data, err := datareader.Read(ctx, "AAPL.US", "stooq", start, end, nil)

// Multiple symbols
symbols := []string{"AAPL.US", "MSFT.US", "^SPX"}
dataMap, err := datareader.Read(ctx, symbols, "stooq", start, end, nil)
```

### Links
- Website: https://stooq.com
- No official API documentation

---

## IEX Cloud

### Overview
IEX Cloud provides professional-grade financial data including real-time and historical stock market data, fundamentals, and more.

### API Key Required
**Yes** - Free tier available (limited)

### Symbol Format
- **US Stocks**: `TICKER` (e.g., `AAPL`, `MSFT`)
- Case-insensitive
- US markets only

### Data Available
- Daily OHLCV data
- Intraday data
- Historical data (up to 5 years on free tier)
- Real-time quotes
- Corporate actions
- Fundamental data

### Rate Limits
- **Free tier**: 50,000 messages/month
- **Paid tiers**: Higher limits
- Each API call consumes messages based on data returned
- Monitor usage in console

### Multi-Symbol Support
✅ **Yes** - Parallel fetching supported

### Capabilities
- ✅ High-quality, exchange-grade data
- ✅ Real-time data available
- ✅ Well-documented API
- ✅ Reliable uptime
- ✅ Rich dataset beyond prices

### Limitations
- ❌ API key required
- ❌ US markets only
- ❌ Free tier has usage limits
- ❌ Message-based pricing

### Example Usage
```go
opts := &datareader.Options{
    APIKey: "your-iex-api-key",
}

data, err := datareader.Read(ctx, "AAPL", "iex", start, end, opts)

// Multiple symbols
symbols := []string{"AAPL", "MSFT", "TSLA"}
dataMap, err := datareader.Read(ctx, symbols, "iex", start, end, opts)
```

### Links
- Website: https://iexcloud.io
- API Documentation: https://iexcloud.io/docs/api
- Console: https://iexcloud.io/console
- Pricing: https://iexcloud.io/pricing

---

## Tiingo

### Overview
Tiingo provides high-quality stock market data with a focus on data accuracy and reliability, including fundamentals and corporate actions.

### API Key Required
**Yes** - Free tier available

### Symbol Format
- **US Stocks**: `TICKER` (e.g., `AAPL`, `MSFT`)
- Case-insensitive

### Data Available
- Daily adjusted OHLCV data
- Historical data (30+ years)
- Fundamental data
- Corporate actions (splits, dividends)
- Intraday data available

### Rate Limits
- **Free tier**: Reasonable limits
- **Paid tiers**: Higher limits
- Check documentation for current limits

### Multi-Symbol Support
⚠️ **Single symbol only** - Call Read() multiple times for multiple symbols

### Capabilities
- ✅ High data quality
- ✅ Dividend and split adjusted data
- ✅ Fundamental data available
- ✅ Long historical coverage
- ✅ Corporate actions included

### Limitations
- ❌ API key required
- ❌ Free tier has limits
- ❌ No parallel multi-symbol fetching yet
- ❌ US markets focus

### Example Usage
```go
opts := &datareader.Options{
    APIKey: "your-tiingo-api-key",
}

data, err := datareader.Read(ctx, "AAPL", "tiingo", start, end, opts)

// For multiple symbols, call separately
for _, symbol := range symbols {
    data, err := datareader.Read(ctx, symbol, "tiingo", start, end, opts)
    // Process each symbol
}
```

### Links
- Website: https://www.tiingo.com
- API Documentation: https://api.tiingo.com/documentation/general/overview
- Get API Key: https://www.tiingo.com/account/api/token

---

## OECD

### Overview
The Organisation for Economic Co-operation and Development (OECD) provides economic indicators and statistics via SDMX-JSON format.

### API Key Required
**No** - Free access without registration

### Symbol Format
- **Dataset Code**: `DATASET` (e.g., `QNA`, `MEI`, `EO`)
- Find dataset codes in OECD data catalog
- Complex multidimensional data structure

### Data Available
- Economic indicators (GDP, inflation, employment, etc.)
- Statistical databases
- Country comparisons
- Time series data
- Historical data (varies by dataset)

### Rate Limits
- No official rate limits
- Recommended: Reasonable request rates
- Large datasets may take time to process

### Multi-Symbol Support
⚠️ **Single dataset only** - Each dataset requires separate call

### Capabilities
- ✅ Free, no API key required
- ✅ Authoritative OECD data
- ✅ Comprehensive economic statistics
- ✅ SDMX standard format
- ✅ International coverage

### Limitations
- ❌ Complex data structure
- ❌ Dataset codes not intuitive
- ❌ SDMX-JSON format requires parsing
- ❌ No parallel multi-dataset fetching

### Example Usage
```go
// Quarterly National Accounts data
data, err := datareader.Read(ctx, "QNA", "oecd", start, end, nil)

// Different datasets require separate calls
datasets := []string{"QNA", "MEI", "EO"}
for _, dataset := range datasets {
    data, err := datareader.Read(ctx, dataset, "oecd", start, end, nil)
    // Process each dataset
}
```

### Links
- Website: https://data.oecd.org
- API Documentation: https://data.oecd.org/api/sdmx-json-documentation/
- Data Catalog: https://stats.oecd.org

---

## Eurostat

### Overview
Eurostat provides official European Union statistics covering economy, population, environment, and more via JSON-stat format.

### API Key Required
**No** - Free access without registration

### Symbol Format
- **Dataset Code**: `DATASET` (e.g., `nama_10_gdp`, `prc_hicp_midx`)
- Find dataset codes on Eurostat website
- Multidimensional data with automatic aggregation

### Data Available
- EU economic indicators
- Population statistics
- Regional data
- Environmental indicators
- Historical data (varies by dataset)

### Rate Limits
- No official rate limits
- Recommended: Reasonable request rates
- Be respectful of server resources

### Multi-Symbol Support
⚠️ **Single dataset only** - Each dataset requires separate call

### Capabilities
- ✅ Free, no API key required
- ✅ Official EU statistics
- ✅ Comprehensive European data
- ✅ JSON-stat format
- ✅ Multidimensional aggregation

### Limitations
- ❌ Complex data structure
- ❌ Dataset codes not intuitive
- ❌ JSON-stat format requires parsing
- ❌ No parallel multi-dataset fetching
- ❌ EU/European focus only

### Example Usage
```go
// GDP dataset
data, err := datareader.Read(ctx, "nama_10_gdp", "eurostat", start, end, nil)

// Price index
data, err := datareader.Read(ctx, "prc_hicp_midx", "eurostat", start, end, nil)
```

### Links
- Website: https://ec.europa.eu/eurostat
- API Documentation: https://ec.europa.eu/eurostat/web/json-and-unicode-web-services
- Data Browser: https://ec.europa.eu/eurostat/databrowser/

---

## Comparison Matrix

| Feature | Yahoo | FRED | World Bank | Alpha Vantage | Stooq | IEX Cloud | Tiingo | OECD | Eurostat |
|---------|-------|------|------------|---------------|-------|-----------|--------|------|----------|
| **API Key Required** | No | Optional | No | Yes | No | Yes | Yes | No | No |
| **Free Tier** | ✅ | ✅ | ✅ | Limited | ✅ | Limited | Limited | ✅ | ✅ |
| **Data Type** | Stocks | Economic | Economic | Stocks | Stocks/Forex | Stocks | Stocks | Economic | Economic |
| **Real-time** | Delayed | No | No | Yes | No | Yes | Yes | No | No |
| **Historical** | 20+ yrs | Varies | 1960+ | 20+ yrs | 10-20 yrs | 5 yrs | 30+ yrs | Varies | Varies |
| **Multi-Symbol** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Rate Limits** | Soft | 120/min | Soft | 5/min | Soft | Message | Moderate | Soft | Soft |
| **Global Markets** | ✅ | ❌ | ✅ | Limited | ✅ | ❌ | ❌ | ✅ | EU Only |
| **Data Quality** | Good | High | High | High | Good | High | High | High | High |
| **Best For** | Stocks | US Econ | Intl Dev | Trading | Free Data | Pro Apps | Quality | OECD Data | EU Stats |

### Legend
- ✅ = Fully supported
- ❌ = Not supported
- Limited = Available with restrictions
- Soft = No strict limits, be reasonable
- Message = Usage-based pricing

---

## Choosing the Right Source

### For Stock Market Data
- **Free & Simple**: Yahoo Finance or Stooq
- **High Quality**: Tiingo or IEX Cloud
- **Real-time Trading**: Alpha Vantage or IEX Cloud
- **International**: Yahoo Finance or Stooq

### For Economic Data
- **US Economy**: FRED
- **International Development**: World Bank
- **OECD Countries**: OECD
- **European Union**: Eurostat

### For Production Applications
- **Reliable & Free**: Yahoo + caching
- **Professional**: IEX Cloud or Tiingo
- **Economic Research**: FRED + World Bank

### For Prototyping
- **No API Key Needed**: Yahoo, Stooq, FRED (optional), World Bank
- **Quick Start**: Yahoo Finance
- **Educational**: Any free source

---

## Rate Limiting Best Practices

1. **Use Caching**: Enable response caching to avoid repeated requests
2. **Respect Limits**: Stay well below stated limits
3. **Implement Backoff**: Use exponential backoff on errors
4. **Monitor Usage**: Track API usage for paid tiers
5. **Batch Requests**: Use multi-symbol fetching where available
6. **Off-Peak Hours**: Schedule large jobs during off-peak times

## API Key Management

```go
// Load from environment variable
opts := &datareader.Options{
    APIKey: os.Getenv("ALPHAVANTAGE_API_KEY"),
}

// Use separate keys for different sources
fredOpts := &datareader.Options{
    APIKey: os.Getenv("FRED_API_KEY"),
}

avOpts := &datareader.Options{
    APIKey: os.Getenv("ALPHAVANTAGE_API_KEY"),
}
```

**Never commit API keys to version control!**

---

For more information, see:
- [Main README](../README.md)
- [API Reference](https://pkg.go.dev/github.com/julianshen/gonp-datareader)
- [Examples](../examples/)
