package controllers

import (
	"net/http"

	"pizza/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InvoiceController struct {
	db *gorm.DB
}

func NewInvoiceController(db *gorm.DB) *InvoiceController {
	return &InvoiceController{db: db}
}

func (c *InvoiceController) GetInvoices(ctx *gin.Context) {
	var invoices []models.Invoice
	if err := c.db.Preload("Items").Find(&invoices).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch invoices"})
		return
	}
	ctx.JSON(http.StatusOK, invoices)
}

func (c *InvoiceController) GetInvoice(ctx *gin.Context) {
	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.db.Preload("Items").First(&invoice, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	ctx.JSON(http.StatusOK, invoice)
}

func (c *InvoiceController) CreateInvoice(ctx *gin.Context) {
	var invoice models.Invoice
	if err := ctx.ShouldBindJSON(&invoice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := c.db.Begin()
	if err := tx.Create(&invoice).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
		return
	}

	// Create invoice items
	for i := range invoice.Items {
		invoice.Items[i].InvoiceID = invoice.ID
		if err := tx.Create(&invoice.Items[i]).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice items"})
			return
		}
	}

	// Commit transaction
	tx.Commit()
	ctx.JSON(http.StatusCreated, invoice)
}

func (c *InvoiceController) UpdateInvoice(ctx *gin.Context) {
	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.db.Preload("Items").First(&invoice, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	var updatedInvoice models.Invoice
	if err := ctx.ShouldBindJSON(&updatedInvoice); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start a transaction
	tx := c.db.Begin()

	// Delete existing items
	if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice items"})
		return
	}

	// Update invoice
	invoice.CustomerName = updatedInvoice.CustomerName
	invoice.CustomerPhone = updatedInvoice.CustomerPhone
	invoice.Date = updatedInvoice.Date
	invoice.Total = updatedInvoice.Total

	if err := tx.Save(&invoice).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice"})
		return
	}

	// Create new items
	for i := range updatedInvoice.Items {
		updatedInvoice.Items[i].InvoiceID = invoice.ID
		if err := tx.Create(&updatedInvoice.Items[i]).Error; err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice items"})
			return
		}
	}

	// Commit transaction
	tx.Commit()
	ctx.JSON(http.StatusOK, invoice)
}

func (c *InvoiceController) DeleteInvoice(ctx *gin.Context) {
	id := ctx.Param("id")
	var invoice models.Invoice
	if err := c.db.First(&invoice, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	// Start a transaction
	tx := c.db.Begin()

	// Delete invoice items
	if err := tx.Where("invoice_id = ?", invoice.ID).Delete(&models.InvoiceItem{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice items"})
		return
	}

	// Delete invoice
	if err := tx.Delete(&invoice).Error; err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invoice"})
		return
	}

	// Commit transaction
	tx.Commit()
	ctx.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
