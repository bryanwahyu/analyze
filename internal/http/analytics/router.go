package analytics

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	appsvc "github.com/bryanwahyu/analyze/internal/application/analyst"
	analyst "github.com/bryanwahyu/analyze/internal/domain/analyst"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Svc *appsvc.Service
}

func NewRouter(svc *appsvc.Service) http.Handler {
	h := &Handler{Svc: svc}
	r := chi.NewRouter()

	r.Get("/analytics/top-genres", h.getTopGenres)
	r.Get("/analytics/top-countries", h.getTopCountries)
	r.Get("/analytics/media-types", h.getMediaTypes)
	r.Get("/analytics/highest-sale", h.getHighestSale)
	r.Get("/analytics/invoice-stats", h.getInvoiceStats)
	r.Get("/analytics/agents", h.getAgents)

	return r
}

// getTopGenres godoc
// @Summary Get top genres by sales
// @Description Returns top N genres by total sales and percentage of total
// @Tags analytics
// @Param limit query int false "Limit the number of genres" default(10)
// @Produce json
// @Success 200 {array} analyst.TopGenre
// @Failure 500 {string} string "internal error"
// @Router /analytics/top-genres [get]
func (h *Handler) getTopGenres(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 10)
	res, err := h.Svc.TopGenres(r.Context(), limit)
	respond(w, res, err)
}

// getTopCountries godoc
// @Summary Get top countries stats
// @Description Returns top N countries by number of customers and average order value
// @Tags analytics
// @Param limit query int false "Limit the number of countries" default(5)
// @Produce json
// @Success 200 {object} analyst.TopCountriesResponse
// @Failure 500 {string} string "internal error"
// @Router /analytics/top-countries [get]
func (h *Handler) getTopCountries(w http.ResponseWriter, r *http.Request) {
	limit := parseIntQuery(r, "limit", 5)
	res, err := h.Svc.TopCountries(r.Context(), limit)
	respond(w, res, err)
}

// getMediaTypes godoc
// @Summary Get media type sales
// @Description Sales by media type with percentage of total
// @Tags analytics
// @Produce json
// @Success 200 {array} analyst.MediaTypeSales
// @Failure 500 {string} string "internal error"
// @Router /analytics/media-types [get]
func (h *Handler) getMediaTypes(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.MediaTypeSales(r.Context())
	respond(w, res, err)
}

// getHighestSale godoc
// @Summary Get highest sale invoices
// @Description Highest invoice total(s) and corresponding customer names
// @Tags analytics
// @Produce json
// @Success 200 {array} analyst.HighestSaleInvoice
// @Failure 500 {string} string "internal error"
// @Router /analytics/highest-sale [get]
func (h *Handler) getHighestSale(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.HighestSaleInvoices(r.Context())
	respond(w, res, err)
}

// getInvoiceStats godoc
// @Summary Get invoice statistics
// @Description Median and mode of invoice totals
// @Tags analytics
// @Produce json
// @Success 200 {object} analyst.InvoiceStats
// @Failure 500 {string} string "internal error"
// @Router /analytics/invoice-stats [get]
func (h *Handler) getInvoiceStats(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.InvoiceStats(r.Context())
	respond(w, res, err)
}

// getAgents godoc
// @Summary Get sales by agents
// @Description Total sales and transactions per support agent
// @Tags analytics
// @Produce json
// @Success 200 {array} analyst.AgentSales
// @Failure 500 {string} string "internal error"
// @Router /analytics/agents [get]
func (h *Handler) getAgents(w http.ResponseWriter, r *http.Request) {
	res, err := h.Svc.SalesByAgents(r.Context())
	respond(w, res, err)
}

func parseIntQuery(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func respond(w http.ResponseWriter, v any, err error) {
	if err != nil {
		log.Printf("analytics: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// keep import of analyst for swagger model references
var _ = analyst.TopGenre{}
