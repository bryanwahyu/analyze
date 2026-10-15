package analyst

import (
	"context"

	domain "github.com/bryanwahyu/analyze/internal/domain/analyst"
)

// Service orchestrates analytics use cases.
type Service struct {
	Repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) TopGenres(ctx context.Context, limit int) ([]domain.TopGenre, error) {
	return s.Repo.TopGenres(ctx, limit)
}

func (s *Service) TopCountries(ctx context.Context, limit int) (domain.TopCountriesResponse, error) {
	customers, err := s.Repo.TopCountriesByCustomers(ctx, limit)
	if err != nil {
		return domain.TopCountriesResponse{}, err
	}
	avgOrder, err := s.Repo.TopCountriesByAvgOrder(ctx, limit)
	if err != nil {
		return domain.TopCountriesResponse{}, err
	}
	return domain.TopCountriesResponse{ByCustomers: customers, ByAvgOrder: avgOrder}, nil
}

func (s *Service) MediaTypeSales(ctx context.Context) ([]domain.MediaTypeSales, error) {
	return s.Repo.MediaTypeSales(ctx)
}

func (s *Service) HighestSaleInvoices(ctx context.Context) ([]domain.HighestSaleInvoice, error) {
	return s.Repo.HighestSaleInvoices(ctx)
}

func (s *Service) InvoiceStats(ctx context.Context) (domain.InvoiceStats, error) {
	return s.Repo.InvoiceStats(ctx)
}

func (s *Service) SalesByAgents(ctx context.Context) ([]domain.AgentSales, error) {
	return s.Repo.SalesByAgents(ctx)
}
