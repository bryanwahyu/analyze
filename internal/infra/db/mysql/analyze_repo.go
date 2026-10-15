package mysql

import (
	"context"
	"database/sql"

	domain "github.com/bryanwahyu/analyze/internal/domain/analyst"
)

type AnalystRepository struct {
	db *sql.DB
}

func NewAnalystRepository(db *sql.DB) *AnalystRepository {
	return &AnalystRepository{db: db}
}

// TopGenres returns top N genres by sales with percentage of total.
func (r *AnalystRepository) TopGenres(ctx context.Context, limit int) ([]domain.TopGenre, error) {
	// Using lower-case snake_case table/column names for Chinook variant.
	// Adjust to your schema if needed.
	const q = `
WITH per_genre AS (
  SELECT g.name AS genre,
         SUM(il.unit_price * il.quantity) AS genre_sales
  FROM invoice_line il
  JOIN track t   ON t.track_id  = il.track_id
  JOIN genre g   ON g.genre_id  = t.genre_id
  GROUP BY g.name
)
SELECT genre,
       genre_sales,
       (genre_sales / NULLIF(SUM(genre_sales) OVER (), 0)) * 100 AS pct_of_total
FROM per_genre
ORDER BY genre_sales DESC
LIMIT ?`

	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.TopGenre{}
	for rows.Next() {
		var e domain.TopGenre
		if err := rows.Scan(&e.Genre, &e.GenreSales, &e.PctOfTotal); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// TopCountriesByCustomers returns top N countries by number of customers.
func (r *AnalystRepository) TopCountriesByCustomers(ctx context.Context, limit int) ([]domain.CountryCustomers, error) {
	const q = `
SELECT country, COUNT(*) AS customer_count
FROM customer
GROUP BY country
ORDER BY customer_count DESC, country
LIMIT ?`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.CountryCustomers{}
	for rows.Next() {
		var e domain.CountryCustomers
		if err := rows.Scan(&e.Country, &e.CustomerCount); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// TopCountriesByAvgOrder returns top N countries by average invoice total.
func (r *AnalystRepository) TopCountriesByAvgOrder(ctx context.Context, limit int) ([]domain.CountryAvgOrder, error) {
	const q = `
SELECT billing_country AS country,
       AVG(total)       AS avg_order_value,
       COUNT(*)         AS orders
FROM invoice
GROUP BY billing_country
ORDER BY avg_order_value DESC, country
LIMIT ?`
	rows, err := r.db.QueryContext(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.CountryAvgOrder{}
	for rows.Next() {
		var e domain.CountryAvgOrder
		if err := rows.Scan(&e.Country, &e.AvgOrderValue, &e.Orders); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// MediaTypeSales returns sales by media type with percentage of total.
func (r *AnalystRepository) MediaTypeSales(ctx context.Context) ([]domain.MediaTypeSales, error) {
	const q = `
WITH per_media AS (
  SELECT mt.name AS media_type,
         SUM(il.unit_price * il.quantity) AS media_sales
  FROM invoice_line il
  JOIN track t       ON t.track_id       = il.track_id
  JOIN media_type mt ON mt.media_type_id = t.media_type_id
  GROUP BY mt.name
)
SELECT media_type,
       media_sales,
       (media_sales / NULLIF(SUM(media_sales) OVER (), 0)) * 100 AS pct_of_total
FROM per_media
ORDER BY media_sales DESC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.MediaTypeSales{}
	for rows.Next() {
		var e domain.MediaTypeSales
		if err := rows.Scan(&e.MediaType, &e.MediaSales, &e.PctOfTotal); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// HighestSaleInvoices returns the invoice(s) with the highest total and their customer names.
func (r *AnalystRepository) HighestSaleInvoices(ctx context.Context) ([]domain.HighestSaleInvoice, error) {
	const q = `
WITH hi AS (
  SELECT MAX(total) AS max_total FROM invoice
)
SELECT i.invoice_id,
       i.total AS highest_total,
       c.first_name,
       c.last_name
FROM invoice i
JOIN hi ON i.total = hi.max_total
JOIN customer c ON c.customer_id = i.customer_id`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.HighestSaleInvoice{}
	for rows.Next() {
		var e domain.HighestSaleInvoice
		if err := rows.Scan(&e.InvoiceID, &e.HighestTotal, &e.FirstName, &e.LastName); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// InvoiceStats returns median and mode of invoice totals.
func (r *AnalystRepository) InvoiceStats(ctx context.Context) (domain.InvoiceStats, error) {
	// Median
	const qMedian = `
WITH ordered AS (
  SELECT total,
         ROW_NUMBER() OVER (ORDER BY total) AS rn,
         COUNT(*) OVER () AS cnt
  FROM invoice
)
SELECT AVG(total) AS median_total
FROM ordered
WHERE rn IN (FLOOR((cnt + 1)/2.0), FLOOR((cnt + 2)/2.0))`

	// Mode
	const qMode = `
SELECT total AS mode_total
FROM invoice
GROUP BY total
ORDER BY COUNT(*) DESC, total DESC
LIMIT 1`

	var stats domain.InvoiceStats

	// median
	if err := r.db.QueryRowContext(ctx, qMedian).Scan(&stats.Median); err != nil {
		return domain.InvoiceStats{}, err
	}
	// mode
	if err := r.db.QueryRowContext(ctx, qMode).Scan(&stats.Mode); err != nil {
		return domain.InvoiceStats{}, err
	}
	return stats, nil
}

// SalesByAgents returns total sales and transactions by support representatives.
func (r *AnalystRepository) SalesByAgents(ctx context.Context) ([]domain.AgentSales, error) {
	const q = `
SELECT e.employee_id,
       CONCAT(e.first_name, ' ', e.last_name) AS agent_name,
       SUM(i.total) AS total_sales,
       COUNT(i.invoice_id) AS transactions
FROM employee e
JOIN customer c ON c.support_rep_id = e.employee_id
JOIN invoice  i ON i.customer_id    = c.customer_id
WHERE e.title = 'Sales Support Agent'
GROUP BY e.employee_id, agent_name
ORDER BY total_sales DESC`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []domain.AgentSales{}
	for rows.Next() {
		var e domain.AgentSales
		if err := rows.Scan(&e.EmployeeID, &e.AgentName, &e.TotalSales, &e.Transactions); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// Close is a no-op because the DB lifecycle is owned by main.
func (r *AnalystRepository) Close() error { return nil }
