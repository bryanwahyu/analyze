package analyst

import "context"

// Repository defines the data-access contract for analytics.
type Repository interface {
    TopGenres(ctx context.Context, limit int) ([]TopGenre, error)
    TopCountriesByCustomers(ctx context.Context, limit int) ([]CountryCustomers, error)
    TopCountriesByAvgOrder(ctx context.Context, limit int) ([]CountryAvgOrder, error)
    MediaTypeSales(ctx context.Context) ([]MediaTypeSales, error)
    HighestSaleInvoices(ctx context.Context) ([]HighestSaleInvoice, error)
    InvoiceStats(ctx context.Context) (InvoiceStats, error)
    SalesByAgents(ctx context.Context) ([]AgentSales, error)
    Close() error
}
