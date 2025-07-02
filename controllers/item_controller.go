package controllers

import (
	"net/http"

	"pizza/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemController struct {
	db *gorm.DB
}

func NewItemController(db *gorm.DB) *ItemController {
	return &ItemController{db: db}
}

func (c *ItemController) GetItems(ctx *gin.Context) {
	var items []models.Item
	if err := c.db.Find(&items).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	ctx.JSON(http.StatusOK, items)
}

func (c *ItemController) GetItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.Item
	if err := c.db.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *ItemController) CreateItem(ctx *gin.Context) {
	var item models.Item
	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.db.Create(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}
	ctx.JSON(http.StatusCreated, item)
}

func (c *ItemController) UpdateItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.Item
	if err := c.db.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&item); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.db.Save(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}
	ctx.JSON(http.StatusOK, item)
}

func (c *ItemController) DeleteItem(ctx *gin.Context) {
	id := ctx.Param("id")
	var item models.Item
	if err := c.db.First(&item, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	if err := c.db.Delete(&item).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete item"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
