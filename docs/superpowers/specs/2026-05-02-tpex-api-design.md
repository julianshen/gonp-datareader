# TPEX API Design

## Scope

Add a `tpex` data source using the official TPEX OpenAPI at `https://www.tpex.org.tw/openapi/`.

This PR supports three independently useful TPEX datasets:

- Mainboard OTC stocks from `/tpex_mainboard_daily_close_quotes`.
- Emerging stocks from `/tpex_esb_latest_statistics`.
- OTC index history from `/tpex_index`.

It does not attempt to wrap every TPEX OpenAPI endpoint. The OpenAPI catalog is broad, and unrelated financial statement, bond, warrant, governance, and ESG endpoints should be added as separate reviewable chunks.

## Public API

The source is registered as `tpex` in `DataReader` and `ListSources`.

Symbol routing:

- `2330` style 4- or 6-digit numeric symbols read mainboard OTC stock quotes.
- `esb:<code>` reads emerging stock quotes.
- `index` reads OTC index history.

The returned type is `*tpex.ParsedData` for `ReadSingle` and `map[string]*tpex.ParsedData` for `Read`, matching the existing TWSE reader shape.

## Data Mapping

Mainboard fields map to OHLCV-style slices:

- `SecuritiesCompanyCode` -> `Symbol`
- `CompanyName` -> `Name`
- `Date` -> `Date`
- `Open`, `High`, `Low`, `Close` -> price slices
- `TradingShares` -> `Volume`
- `TransactionNumber` -> `Transactions`
- `Change` -> `Change`

Emerging stock fields map as:

- `LatestPrice` -> `Close`
- `Highest` -> `High`
- `Lowest` -> `Low`
- `Average` -> `Open`
- `TransactionVolume` -> `Volume`

Index fields map as:

- `Open`, `High`, `Low`, `Close`, `Change` -> matching price slices
- volume and transaction slices remain zero-valued.

## Error Handling

The reader validates symbols and date ranges before HTTP requests. HTTP non-200 responses return status errors. Missing symbols return explicit not-found errors. Numeric parsing accepts commas and leading `+` signs because TPEX publishes localized numeric strings.

## Testing

Follow strict TDD:

- Parser tests for mainboard, emerging stock, and index schemas.
- Reader tests with `httptest` servers to verify endpoint routing and filtering.
- Factory tests for `DataReader("tpex")` and `ListSources`.
- Coverage must remain above 90% using the repository coverage command.
