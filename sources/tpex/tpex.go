// Package tpex provides data access to Taipei Exchange (TPEx).
package tpex

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/internal/utils"
	"github.com/julianshen/gonp-datareader/sources"
)

const (
	tpexBaseURL       = "https://www.tpex.org.tw/openapi"
	mainboardEndpoint = "/tpex_mainboard_daily_close_quotes"
	emergingEndpoint  = "/tpex_esb_latest_statistics"
	indexEndpoint     = "/tpex_index"
)

var tpexNumericSymbolPattern = regexp.MustCompile(`^[0-9]{4}$|^[0-9]{6}$`)

// TPEXReader fetches data from Taipei Exchange.
type TPEXReader struct {
	*sources.BaseSource
	client  *internalhttp.RetryableClient
	baseURL string
}

// NewTPEXReader creates a TPEX reader with the default OpenAPI base URL.
func NewTPEXReader(opts *internalhttp.ClientOptions) *TPEXReader {
	return NewTPEXReaderWithBaseURL(opts, tpexBaseURL)
}

// NewTPEXReaderWithBaseURL creates a TPEX reader with a custom base URL for tests.
func NewTPEXReaderWithBaseURL(opts *internalhttp.ClientOptions, baseURL string) *TPEXReader {
	if opts == nil {
		opts = internalhttp.DefaultClientOptions()
	}
	return &TPEXReader{
		BaseSource: sources.NewBaseSource("tpex"),
		client:     internalhttp.NewRetryableClient(opts),
		baseURL:    baseURL,
	}
}

// Name returns the display name of the data source.
func (t *TPEXReader) Name() string {
	return "Taipei Exchange"
}

// ValidateSymbol checks whether a TPEX symbol is routable.
func (t *TPEXReader) ValidateSymbol(symbol string) error {
	if symbol == "" {
		return fmt.Errorf("symbol cannot be empty")
	}
	if strings.ContainsAny(symbol, " \t\n\r") {
		return fmt.Errorf("symbol contains invalid characters")
	}
	if symbol == "index" {
		return nil
	}
	if strings.HasPrefix(symbol, "esb:") {
		code := strings.TrimPrefix(symbol, "esb:")
		if tpexNumericSymbolPattern.MatchString(code) {
			return nil
		}
		return fmt.Errorf("invalid TPEX emerging stock symbol format: %q (must be esb:<4 or 6 digits>)", symbol)
	}
	if tpexNumericSymbolPattern.MatchString(symbol) {
		return nil
	}
	return fmt.Errorf("invalid TPEX symbol format: %q", symbol)
}

// BuildURL returns the default mainboard endpoint URL.
func (t *TPEXReader) BuildURL() string {
	return buildEndpointURL(t.baseURL, mainboardEndpoint)
}

func buildEndpointURL(baseURL, endpoint string) string {
	return strings.TrimRight(baseURL, "/") + endpoint
}

// ReadSingle fetches one TPEX symbol.
func (t *TPEXReader) ReadSingle(ctx context.Context, symbol string, start, end time.Time) (interface{}, error) {
	if err := t.ValidateSymbol(symbol); err != nil {
		return nil, fmt.Errorf("invalid symbol: %w", err)
	}
	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	data, err := t.readBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return filterByDateRange(data, start, end), nil
}

// Read fetches multiple TPEX symbols.
func (t *TPEXReader) Read(ctx context.Context, symbols []string, start, end time.Time) (interface{}, error) {
	if len(symbols) == 0 {
		return nil, fmt.Errorf("invalid symbols: symbol list cannot be empty")
	}
	for _, symbol := range symbols {
		if err := t.ValidateSymbol(symbol); err != nil {
			return nil, fmt.Errorf("invalid symbols: %w", err)
		}
	}
	if err := utils.ValidateDateRange(start, end); err != nil {
		return nil, fmt.Errorf("invalid date range: %w", err)
	}

	result := make(map[string]*ParsedData, len(symbols))
	for _, symbol := range symbols {
		data, err := t.ReadSingle(ctx, symbol, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", symbol, err)
		}
		result[symbol] = data.(*ParsedData)
	}
	return result, nil
}

func (t *TPEXReader) readBySymbol(ctx context.Context, symbol string) (*ParsedData, error) {
	switch {
	case symbol == "index":
		return t.readIndex(ctx)
	case strings.HasPrefix(symbol, "esb:"):
		return t.readEmerging(ctx, symbol)
	default:
		return t.readMainboard(ctx, symbol)
	}
}

func (t *TPEXReader) readMainboard(ctx context.Context, symbol string) (*ParsedData, error) {
	body, err := t.get(ctx, mainboardEndpoint)
	if err != nil {
		return nil, err
	}
	rows, err := parseMainboardJSON(body)
	if err != nil {
		return nil, err
	}
	row, err := filterMainboardBySymbol(rows, symbol)
	if err != nil {
		return nil, err
	}
	return parseMainboardData(row)
}

func (t *TPEXReader) readEmerging(ctx context.Context, symbol string) (*ParsedData, error) {
	body, err := t.get(ctx, emergingEndpoint)
	if err != nil {
		return nil, err
	}
	rows, err := parseEmergingJSON(body)
	if err != nil {
		return nil, err
	}
	row, err := filterEmergingBySymbol(rows, symbol)
	if err != nil {
		return nil, err
	}
	return parseEmergingData(row)
}

func (t *TPEXReader) readIndex(ctx context.Context) (*ParsedData, error) {
	body, err := t.get(ctx, indexEndpoint)
	if err != nil {
		return nil, err
	}
	rows, err := parseIndexJSON(body)
	if err != nil {
		return nil, err
	}
	combined := &ParsedData{Symbol: "index", Name: "TPEx Index"}
	for _, row := range rows {
		data, err := parseIndexData(row)
		if err != nil {
			return nil, err
		}
		appendParsedData(combined, data)
	}
	return combined, nil
}

func (t *TPEXReader) get(ctx context.Context, endpoint string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, buildEndpointURL(t.baseURL, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := t.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch data: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	return body, nil
}

func appendParsedData(dst, src *ParsedData) {
	dst.Date = append(dst.Date, src.Date...)
	dst.Open = append(dst.Open, src.Open...)
	dst.High = append(dst.High, src.High...)
	dst.Low = append(dst.Low, src.Low...)
	dst.Close = append(dst.Close, src.Close...)
	dst.Volume = append(dst.Volume, src.Volume...)
	dst.Transactions = append(dst.Transactions, src.Transactions...)
	dst.Change = append(dst.Change, src.Change...)
}
