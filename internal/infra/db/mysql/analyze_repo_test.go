package mysql

import (
	"context"
	"math"
	"os"
	"testing"
)

// These tests run against the Chinook script in DB.sql. Start it with
// `docker compose up -d db` and set
// MYSQL_DSN="root:example@tcp(127.0.0.1:3306)/music_analyze?parseTime=true".
func testRepo(t *testing.T) *AnalystRepository {
	t.Helper()
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		t.Skip("MYSQL_DSN not set")
	}
	db, err := Connect(context.Background(), dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewAnalystRepository(db)
}

func near(a, b float64) bool { return math.Abs(a-b) < 0.01 }

func TestTopGenres(t *testing.T) {
	got, err := testRepo(t).TopGenres(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 || got[0].Genre != "Rock" || !near(got[0].GenreSales, 826.65) {
		t.Fatalf("top genres = %+v", got)
	}
	if !(got[0].GenreSales >= got[1].GenreSales && got[1].GenreSales >= got[2].GenreSales) {
		t.Errorf("not sorted by sales: %+v", got)
	}
}

func TestMediaTypeSharesAddUpTo100(t *testing.T) {
	got, err := testRepo(t).MediaTypeSales(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var sum float64
	for _, m := range got {
		sum += m.PctOfTotal
	}
	if math.Abs(sum-100) > 0.01 {
		t.Errorf("percentages sum to %.4f", sum)
	}
}

func TestTopCountries(t *testing.T) {
	r := testRepo(t)
	byCustomers, err := r.TopCountriesByCustomers(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(byCustomers) != 1 || byCustomers[0].Country != "USA" || byCustomers[0].CustomerCount != 13 {
		t.Errorf("by customers = %+v", byCustomers)
	}
	byAvg, err := r.TopCountriesByAvgOrder(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(byAvg) != 1 || byAvg[0].Country != "Chile" || !near(byAvg[0].AvgOrderValue, 6.66) {
		t.Errorf("by avg order = %+v", byAvg)
	}
}

func TestHighestSaleInvoice(t *testing.T) {
	got, err := testRepo(t).HighestSaleInvoices(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].InvoiceID != 404 || !near(got[0].HighestTotal, 25.86) || got[0].LastName != "Holý" {
		t.Errorf("highest sale = %+v", got)
	}
}

func TestInvoiceStats(t *testing.T) {
	got, err := testRepo(t).InvoiceStats(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !near(got.Median, 3.96) || !near(got.Mode, 1.98) {
		t.Errorf("stats = %+v", got)
	}
}

func TestSalesByAgents(t *testing.T) {
	got, err := testRepo(t).SalesByAgents(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	var transactions int
	for _, a := range got {
		transactions += a.Transactions
	}
	if len(got) != 3 || got[0].AgentName != "Jane Peacock" || transactions != 412 {
		t.Errorf("agents = %+v (transactions %d)", got, transactions)
	}
}

func TestEmptyResultIsEmptySliceNotNil(t *testing.T) {
	got, err := testRepo(t).TopGenres(context.Background(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("want empty slice so the API returns [] instead of null")
	}
}
