package main

import (
	"log"

	"pizza/backend/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := "host=localhost user=postgres password=postgres dbname=pizza_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// Initialize controllers
	setupRoutes(r, db)

	// Start server
	log.Fatal(r.Run(":8080"))
}

func setupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize controllers
	itemController := controllers.NewItemController(db)
	invoiceController := controllers.NewInvoiceController(db)
	dashboardController := controllers.NewDashboardController(db)

	// Item routes
	r.GET("/items", itemController.GetItems)
	r.GET("/items/:id", itemController.GetItem)
	r.POST("/items", itemController.CreateItem)
	r.PUT("/items/:id", itemController.UpdateItem)
	r.DELETE("/items/:id", itemController.DeleteItem)

	// Invoice routes
	r.GET("/invoices", invoiceController.GetInvoices)
	r.GET("/invoices/:id", invoiceController.GetInvoice)
	r.POST("/invoices", invoiceController.CreateInvoice)
	r.PUT("/invoices/:id", invoiceController.UpdateInvoice)
	r.DELETE("/invoices/:id", invoiceController.DeleteInvoice)

	// Dashboard routes
	r.GET("/dashboard/metrics", dashboardController.GetDashboardMetrics)
	r.GET("/dashboard/top-items", dashboardController.GetTopSellingItems)
	r.GET("/dashboard/recent-orders", dashboardController.GetRecentOrders)
}
