package analytics

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	appsvc "github.com/bryanwahyu/analyze/internal/application/analyst"
	domain "github.com/bryanwahyu/analyze/internal/domain/analyst"
)

type fakeRepo struct {
	gotLimit int
	err      error
}

func (f *fakeRepo) TopGenres(_ context.Context, limit int) ([]domain.TopGenre, error) {
	f.gotLimit = limit
	return []domain.TopGenre{{Genre: "Rock", GenreSales: 826.65, PctOfTotal: 35.5}}, f.err
}
func (f *fakeRepo) TopCountriesByCustomers(_ context.Context, limit int) ([]domain.CountryCustomers, error) {
	f.gotLimit = limit
	return []domain.CountryCustomers{{Country: "USA", CustomerCount: 13}}, f.err
}
func (f *fakeRepo) TopCountriesByAvgOrder(context.Context, int) ([]domain.CountryAvgOrder, error) {
	return []domain.CountryAvgOrder{{Country: "Chile", AvgOrderValue: 6.66, Orders: 7}}, f.err
}
func (f *fakeRepo) MediaTypeSales(context.Context) ([]domain.MediaTypeSales, error) {
	return nil, f.err
}
func (f *fakeRepo) HighestSaleInvoices(context.Context) ([]domain.HighestSaleInvoice, error) {
	return nil, f.err
}
func (f *fakeRepo) InvoiceStats(context.Context) (domain.InvoiceStats, error) {
	return domain.InvoiceStats{Median: 3.96, Mode: 1.98}, f.err
}
func (f *fakeRepo) SalesByAgents(context.Context) ([]domain.AgentSales, error) { return nil, f.err }
func (f *fakeRepo) Close() error                                               { return nil }

func get(t *testing.T, repo *fakeRepo, target string) *httptest.ResponseRecorder {
	t.Helper()
	rec := httptest.NewRecorder()
	NewRouter(appsvc.NewService(repo)).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
	return rec
}

func TestTopGenres_LimitQuery(t *testing.T) {
	cases := map[string]int{
		"/analytics/top-genres":          10,
		"/analytics/top-genres?limit=3":  3,
		"/analytics/top-genres?limit=0":  10,
		"/analytics/top-genres?limit=-2": 10,
		"/analytics/top-genres?limit=ab": 10,
	}
	for target, want := range cases {
		repo := &fakeRepo{}
		rec := get(t, repo, target)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s: status %d", target, rec.Code)
		}
		if repo.gotLimit != want {
			t.Errorf("%s: limit = %d, want %d", target, repo.gotLimit, want)
		}
	}
}

func TestTopCountries_CombinesBothRankings(t *testing.T) {
	rec := get(t, &fakeRepo{}, "/analytics/top-countries")

	var body domain.TopCountriesResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.ByCustomers) != 1 || body.ByCustomers[0].Country != "USA" {
		t.Errorf("by_customers = %+v", body.ByCustomers)
	}
	if len(body.ByAvgOrder) != 1 || body.ByAvgOrder[0].Country != "Chile" {
		t.Errorf("by_avg_order = %+v", body.ByAvgOrder)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("content-type = %q", ct)
	}
}

func TestRepositoryErrorIsNotLeaked(t *testing.T) {
	rec := get(t, &fakeRepo{err: errors.New("Error 1146: Table 'music_analyze.invoice' doesn't exist")}, "/analytics/invoice-stats")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "music_analyze") {
		t.Errorf("response leaks the database error: %q", rec.Body.String())
	}
}
