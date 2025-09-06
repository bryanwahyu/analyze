package analyst

// Core analytics entities used across repository, service, and HTTP layers.

type TopGenre struct {
    Genre       string  `json:"genre"`
    GenreSales  float64 `json:"genre_sales"`
    PctOfTotal  float64 `json:"pct_of_total"`
}

type CountryCustomers struct {
    Country       string `json:"country"`
    CustomerCount int    `json:"customer_count"`
}

type CountryAvgOrder struct {
    Country       string  `json:"country"`
    AvgOrderValue float64 `json:"avg_order_value"`
    Orders        int     `json:"orders"`
}

type MediaTypeSales struct {
    MediaType  string  `json:"media_type"`
    MediaSales float64 `json:"media_sales"`
    PctOfTotal float64 `json:"pct_of_total"`
}

type HighestSaleInvoice struct {
    InvoiceID   int     `json:"invoice_id"`
    HighestTotal float64 `json:"highest_total"`
    FirstName   string  `json:"first_name"`
    LastName    string  `json:"last_name"`
}

type InvoiceStats struct {
    Median float64 `json:"median_total"`
    Mode   float64 `json:"mode_total"`
}

type AgentSales struct {
    EmployeeID  int     `json:"employee_id"`
    AgentName   string  `json:"agent_name"`
    TotalSales  float64 `json:"total_sales"`
    Transactions int    `json:"transactions"`
}

// Aggregated response for the Top Countries (customers and avg order) use case.
type TopCountriesResponse struct {
    ByCustomers []CountryCustomers `json:"by_customers"`
    ByAvgOrder  []CountryAvgOrder  `json:"by_avg_order"`
}
