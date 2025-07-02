package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"pizza/backend/models"
)

type DashboardController struct {
	db *gorm.DB
}

type DashboardMetrics struct {
	TotalRevenue     float64 `json:"totalRevenue"`
	TotalOrders      int     `json:"totalOrders"`
	AverageOrderValue float64 `json:"averageOrderValue"`
	TotalItems       int     `json:"totalItems"`
	RecentRevenue    float64 `json:"recentRevenue"`
	RecentOrders     int     `json:"recentOrders"`
	GrowthRate       float64 `json:"growthRate"`
}

type TopSellingItem struct {
	ItemName string  `json:"itemName"`
	Quantity int     `json:"quantity"`
	Revenue  float64 `json:"revenue"`
}

type RecentOrder struct {
	InvoiceNumber string    `json:"invoiceNumber"`
	CustomerName  string    `json:"customerName"`
	ItemCount     int       `json:"itemCount"`
	Total        float64   `json:"total"`
	Date         time.Time `json:"date"`
}

func NewDashboardController(db *gorm.DB) *DashboardController {
	return &DashboardController{db: db}
}

func (dc *DashboardController) GetDashboardMetrics(c *gin.Context) {
	var metrics DashboardMetrics

	// Calculate total revenue and orders
	var invoices []models.Invoice
	if err := dc.db.Find(&invoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}

	// Count total items
	var totalItems int64
	if err := dc.db.Model(&models.Item{}).Count(&totalItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count items"})
		return
	}

	// Calculate metrics
	for _, invoice := range invoices {
		metrics.TotalRevenue += invoice.Total
	}
	metrics.TotalOrders = len(invoices)
	if metrics.TotalOrders > 0 {
		metrics.AverageOrderValue = metrics.TotalRevenue / float64(metrics.TotalOrders)
	}
	metrics.TotalItems = int(totalItems)

	// Calculate recent metrics (last 7 days)
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var recentInvoices []models.Invoice
	if err := dc.db.Where("date >= ?", sevenDaysAgo).Find(&recentInvoices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent invoices"})
		return
	}

	for _, invoice := range recentInvoices {
		metrics.RecentRevenue += invoice.Total
	}
	metrics.RecentOrders = len(recentInvoices)

	// Calculate growth rate (comparing current week with previous week)
	previousWeekStart := sevenDaysAgo.AddDate(0, 0, -7)
	var previousWeekRevenue float64
	dc.db.Model(&models.Invoice{}).
		Where("date >= ? AND date < ?", previousWeekStart, sevenDaysAgo).
		Select("COALESCE(SUM(total), 0)").
		Row().
		Scan(&previousWeekRevenue)

	if previousWeekRevenue > 0 {
		metrics.GrowthRate = ((metrics.RecentRevenue - previousWeekRevenue) / previousWeekRevenue) * 100
	}

	c.JSON(http.StatusOK, metrics)
}

func (dc *DashboardController) GetTopSellingItems(c *gin.Context) {
	var topItems []struct {
		ItemName string
		Quantity int
		Revenue  float64
	}

	// Join invoices and items to get sales data
	query := `
		SELECT 
			i.name as item_name,
			SUM(ii.quantity) as quantity,
			SUM(ii.quantity * ii.price) as revenue
		FROM items i
		JOIN invoice_items ii ON i.id = ii.item_id
		GROUP BY i.id, i.name
		ORDER BY quantity DESC
		LIMIT 5
	`

	if err := dc.db.Raw(query).Scan(&topItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch top selling items"})
		return
	}

	result := make([]TopSellingItem, len(topItems))
	for i, item := range topItems {
		result[i] = TopSellingItem{
			ItemName: item.ItemName,
			Quantity: item.Quantity,
			Revenue:  item.Revenue,
		}
	}

	c.JSON(http.StatusOK, result)
}

func (dc *DashboardController) GetRecentOrders(c *gin.Context) {
	var recentOrders []struct {
		InvoiceNumber string
		CustomerName  string
		ItemCount     int
		Total         float64
		Date          time.Time
	}

	query := `
		SELECT 
			i.invoice_number,
			i.customer_name,
			COUNT(ii.item_id) as item_count,
			i.total,
			i.date
		FROM invoices i
		LEFT JOIN invoice_items ii ON i.id = ii.invoice_id
		GROUP BY i.id, i.invoice_number, i.customer_name, i.total, i.date
		ORDER BY i.date DESC
		LIMIT 5
	`

	if err := dc.db.Raw(query).Scan(&recentOrders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent orders"})
		return
	}

	result := make([]RecentOrder, len(recentOrders))
	for i, order := range recentOrders {
		result[i] = RecentOrder{
			InvoiceNumber: order.InvoiceNumber,
			CustomerName:  order.CustomerName,
			ItemCount:     order.ItemCount,
			Total:         order.Total,
			Date:         order.Date,
		}
	}

	c.JSON(http.StatusOK, result)
} 