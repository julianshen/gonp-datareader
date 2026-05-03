package tpex

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/julianshen/gonp-datareader/internal/http"
	"github.com/julianshen/gonp-datareader/sources"
)

func TestTPEXReader_ImplementsReader(t *testing.T) {
	reader := NewTPEXReader(nil)
	var _ sources.Reader = reader

	if reader.Name() != "Taipei Exchange" {
		t.Fatalf("Name() = %q, want Taipei Exchange", reader.Name())
	}
	if reader.Source() != "tpex" {
		t.Fatalf("Source() = %q, want tpex", reader.Source())
	}
}

func TestTPEXReader_ValidateSymbol(t *testing.T) {
	reader := NewTPEXReader(nil)
	tests := []struct {
		symbol  string
		wantErr bool
	}{
		{symbol: "8069"},
		{symbol: "00679B"},
		{symbol: "esb:6871"},
		{symbol: "index"},
		{symbol: "", wantErr: true},
		{symbol: "AAPL.US", wantErr: true},
		{symbol: "esb:", wantErr: true},
	}

	for _, tt := range tests {
		name := tt.symbol
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			err := reader.ValidateSymbol(tt.symbol)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateSymbol(%q) error = %v, wantErr %v", tt.symbol, err, tt.wantErr)
			}
		})
	}
}

func TestBuildEndpointURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		endpoint string
		want     string
	}{
		{
			name:     "without trailing slash",
			baseURL:  "https://www.tpex.org.tw/openapi/v1",
			endpoint: mainboardEndpoint,
			want:     "https://www.tpex.org.tw/openapi/v1/tpex_mainboard_daily_close_quotes",
		},
		{
			name:     "with trailing slash",
			baseURL:  "https://www.tpex.org.tw/openapi/v1/",
			endpoint: indexEndpoint,
			want:     "https://www.tpex.org.tw/openapi/v1/tpex_index",
		},
	}

	for _, tt := range tests {
		got := buildEndpointURL(tt.baseURL, tt.endpoint)
		if got != tt.want {
			t.Fatalf("buildEndpointURL() = %q, want %q", got, tt.want)
		}
	}
}

func TestTPEXReader_BuildURL(t *testing.T) {
	reader := NewTPEXReader(nil)
	if got := reader.BuildURL(); got != "https://www.tpex.org.tw/openapi/v1/tpex_mainboard_daily_close_quotes" {
		t.Fatalf("BuildURL() = %q", got)
	}
}

func TestTPEXReader_ReadSingle_RoutesMainboard(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != mainboardEndpoint {
			t.Fatalf("path = %q, want %q", r.URL.Path, mainboardEndpoint)
		}
		writeJSON(t, w, []tpexMainboardQuote{{
			Date:                  "2025/10/31",
			SecuritiesCompanyCode: "8069",
			CompanyName:           "元太",
			Open:                  "212.50",
			High:                  "218.00",
			Low:                   "210.00",
			Close:                 "216.50",
			TradingShares:         "12,345",
			TransactionNumber:     "67",
		}})
	})

	got := mustReadSingle(t, reader, "8069")
	if got.Symbol != "8069" || got.Close[0] != 216.50 {
		t.Fatalf("got %+v", got)
	}
}

func TestTPEXReader_ReadSingle_RoutesEmerging(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != emergingEndpoint {
			t.Fatalf("path = %q, want %q", r.URL.Path, emergingEndpoint)
		}
		writeJSON(t, w, []tpexEmergingQuote{{
			Date:                  "2025/10/31",
			SecuritiesCompanyCode: "6871",
			CompanyName:           "訊芯",
			Average:               "95.10",
			Highest:               "98.20",
			Lowest:                "93.00",
			LatestPrice:           "97.50",
			TransactionVolume:     "1,250",
		}})
	})

	got := mustReadSingle(t, reader, "esb:6871")
	if got.Symbol != "esb:6871" || got.Close[0] != 97.50 {
		t.Fatalf("got %+v", got)
	}
}

func TestTPEXReader_ReadSingle_RoutesIndex(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != indexEndpoint {
			t.Fatalf("path = %q, want %q", r.URL.Path, indexEndpoint)
		}
		writeJSON(t, w, []tpexIndexData{
			{Date: "2025/10/30", Open: "250.10", High: "252.20", Low: "249.30", Close: "251.80"},
			{Date: "2025/10/31", Open: "251.80", High: "253.00", Low: "250.50", Close: "252.40"},
		})
	})

	got := mustReadSingle(t, reader, "index")
	if got.Symbol != "index" || len(got.Date) != 2 || got.Close[1] != 252.40 {
		t.Fatalf("got %+v", got)
	}
}

func TestTPEXReader_ReadSingle_NotFound(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, []tpexMainboardQuote{})
	})

	_, err := reader.ReadSingle(context.Background(), "8069", testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("ReadSingle() error = %v, want not found", err)
	}
}

func TestTPEXReader_ReadSingle_IndexNoData(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != indexEndpoint {
			t.Fatalf("path = %q, want %q", r.URL.Path, indexEndpoint)
		}
		writeJSON(t, w, []tpexIndexData{})
	})

	_, err := reader.ReadSingle(context.Background(), "index", testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "no index data") {
		t.Fatalf("ReadSingle() error = %v, want no index data", err)
	}
}

func TestTPEXReader_ReadSingle_FiltersIndexDateRange(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, []tpexIndexData{
			{Date: "2025/09/30", Open: "1", High: "2", Low: "3", Close: "4"},
			{Date: "2025/10/31", Open: "5", High: "6", Low: "7", Close: "8"},
		})
	})

	got := mustReadSingle(t, reader, "index")
	if len(got.Date) != 1 || got.Close[0] != 8 {
		t.Fatalf("filtered index data = %+v, want one in-range row", got)
	}
}

func TestTPEXReader_ReadSingle_ParseErrors(t *testing.T) {
	tests := []struct {
		name   string
		symbol string
		path   string
		body   string
	}{
		{
			name:   "mainboard invalid json",
			symbol: "8069",
			path:   mainboardEndpoint,
			body:   "bad",
		},
		{
			name:   "emerging invalid json",
			symbol: "esb:6871",
			path:   emergingEndpoint,
			body:   "bad",
		},
		{
			name:   "index invalid json",
			symbol: "index",
			path:   indexEndpoint,
			body:   "bad",
		},
		{
			name:   "index invalid row",
			symbol: "index",
			path:   indexEndpoint,
			body:   `[{"Date":"bad"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Fatalf("path = %q, want %q", r.URL.Path, tt.path)
				}
				writeString(t, w, tt.body)
			})

			_, err := reader.ReadSingle(context.Background(), tt.symbol, testStart(), testEnd())
			if err == nil {
				t.Fatal("ReadSingle() error = nil, want error")
			}
		})
	}
}

func TestTPEXReader_ReadSingle_HTTPStatus(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	})

	_, err := reader.ReadSingle(context.Background(), "8069", testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "HTTP 500") {
		t.Fatalf("ReadSingle() error = %v, want HTTP 500", err)
	}
}

func TestTPEXReader_ReadSingle_ValidationErrors(t *testing.T) {
	reader := NewTPEXReader(nil)
	_, err := reader.ReadSingle(context.Background(), "bad", testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "invalid symbol") {
		t.Fatalf("ReadSingle invalid symbol error = %v", err)
	}

	_, err = reader.ReadSingle(context.Background(), "8069", testEnd(), testStart())
	if err == nil || !strings.Contains(err.Error(), "invalid date range") {
		t.Fatalf("ReadSingle invalid date error = %v", err)
	}
}

func TestTPEXReader_Read(t *testing.T) {
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case mainboardEndpoint:
			writeJSON(t, w, []tpexMainboardQuote{{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "8069",
				CompanyName:           "元太",
				Open:                  "212.50",
				High:                  "218.00",
				Low:                   "210.00",
				Close:                 "216.50",
			}})
		case emergingEndpoint:
			writeJSON(t, w, []tpexEmergingQuote{{
				Date:                  "2025/10/31",
				SecuritiesCompanyCode: "6871",
				CompanyName:           "訊芯",
				Average:               "95.10",
				Highest:               "98.20",
				Lowest:                "93.00",
				LatestPrice:           "97.50",
			}})
		case indexEndpoint:
			writeJSON(t, w, []tpexIndexData{{Date: "2025/10/31", Open: "1", High: "2", Low: "3", Close: "4"}})
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	result, err := reader.Read(context.Background(), []string{"8069", "esb:6871", "index"}, testStart(), testEnd())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	data := result.(map[string]*ParsedData)
	if len(data) != 3 || data["8069"].Close[0] != 216.50 || data["esb:6871"].Close[0] != 97.50 || data["index"].Close[0] != 4 {
		t.Fatalf("Read() = %+v", data)
	}
}

func TestTPEXReader_Read_FetchesEachEndpointOnce(t *testing.T) {
	counts := map[string]int{}
	reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
		counts[r.URL.Path]++
		switch r.URL.Path {
		case mainboardEndpoint:
			writeJSON(t, w, []tpexMainboardQuote{
				{Date: "2025/10/31", SecuritiesCompanyCode: "8069", CompanyName: "元太", Open: "1", High: "2", Low: "3", Close: "4"},
				{Date: "2025/10/31", SecuritiesCompanyCode: "00679B", CompanyName: "元大美債20年", Open: "5", High: "6", Low: "7", Close: "8"},
			})
		case emergingEndpoint:
			writeJSON(t, w, []tpexEmergingQuote{{Date: "2025/10/31", SecuritiesCompanyCode: "6871", CompanyName: "訊芯", Highest: "2", Lowest: "1", LatestPrice: "1.5"}})
		case indexEndpoint:
			writeJSON(t, w, []tpexIndexData{{Date: "2025/10/31", Open: "1", High: "2", Low: "3", Close: "4"}})
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	})

	_, err := reader.Read(context.Background(), []string{"8069", "00679B", "esb:6871", "index"}, testStart(), testEnd())
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if counts[mainboardEndpoint] != 1 || counts[emergingEndpoint] != 1 || counts[indexEndpoint] != 1 {
		t.Fatalf("endpoint counts = %+v, want one fetch per endpoint", counts)
	}
}

func TestTPEXReader_Read_ValidationErrors(t *testing.T) {
	reader := NewTPEXReader(nil)
	_, err := reader.Read(context.Background(), nil, testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "symbol list cannot be empty") {
		t.Fatalf("Read nil symbols error = %v", err)
	}

	_, err = reader.Read(context.Background(), []string{"bad"}, testStart(), testEnd())
	if err == nil || !strings.Contains(err.Error(), "invalid symbols") {
		t.Fatalf("Read invalid symbol error = %v", err)
	}

	_, err = reader.Read(context.Background(), []string{"8069"}, testEnd(), testStart())
	if err == nil || !strings.Contains(err.Error(), "invalid date range") {
		t.Fatalf("Read invalid date error = %v", err)
	}
}

func TestTPEXReader_Read_PropagatesGroupedReadErrors(t *testing.T) {
	tests := []struct {
		name    string
		symbols []string
		path    string
		body    string
		want    string
	}{
		{
			name:    "emerging not found",
			symbols: []string{"esb:9999"},
			path:    emergingEndpoint,
			body:    `[{"Date":"2025/10/31","SecuritiesCompanyCode":"6871","Highest":"2","Lowest":"1","LatestPrice":"1.5"}]`,
			want:    "not found",
		},
		{
			name:    "index no data",
			symbols: []string{"index"},
			path:    indexEndpoint,
			body:    `[]`,
			want:    "no index data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := newTestReader(t, func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != tt.path {
					t.Fatalf("path = %q, want %q", r.URL.Path, tt.path)
				}
				writeString(t, w, tt.body)
			})

			_, err := reader.Read(context.Background(), tt.symbols, testStart(), testEnd())
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Read() error = %v, want %q", err, tt.want)
			}
		})
	}
}

func newTestReader(t *testing.T, handler http.HandlerFunc) *TPEXReader {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	opts := internalhttp.DefaultClientOptions()
	return NewTPEXReaderWithBaseURL(opts, server.URL)
}

func writeJSON(t *testing.T, w http.ResponseWriter, value interface{}) {
	t.Helper()
	if err := json.NewEncoder(w).Encode(value); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}

func writeString(t *testing.T, w http.ResponseWriter, value string) {
	t.Helper()
	if _, err := w.Write([]byte(value)); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func mustReadSingle(t *testing.T, reader *TPEXReader, symbol string) *ParsedData {
	t.Helper()
	result, err := reader.ReadSingle(context.Background(), symbol, testStart(), testEnd())
	if err != nil {
		t.Fatalf("ReadSingle(%q) error = %v", symbol, err)
	}
	data, ok := result.(*ParsedData)
	if !ok {
		t.Fatalf("ReadSingle(%q) returned %T, want *ParsedData", symbol, result)
	}
	return data
}

func testStart() time.Time {
	return time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
}

func testEnd() time.Time {
	return time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)
}
