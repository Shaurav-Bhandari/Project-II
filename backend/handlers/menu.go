package handlers

import (
	"net/http"
	"sort"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ========== MENU HANDLERS ==========

func (h *Handler) GetCategories(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, name, description, display_order, is_active, created_at
		FROM menu_categories
		ORDER BY display_order, name
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var categories []models.MenuCategory
	for rows.Next() {
		var cat models.MenuCategory
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.DisplayOrder, &cat.IsActive, &cat.CreatedAt); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, categories)
}

type CreateCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var catID uuid.UUID
	err := h.db.QueryRow(`
		INSERT INTO menu_categories (name, description, display_order)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.Name, req.Description, req.DisplayOrder).Scan(&catID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": catID, "message": "Category created"})
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE menu_categories SET name = $1, description = $2, display_order = $3
		WHERE id = $4
	`, req.Name, req.Description, req.DisplayOrder, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category updated"})
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE menu_categories SET is_active = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

func (h *Handler) GetMenuItems(c *gin.Context) {
	categoryID := c.Query("category_id")

	query := `
		SELECT id, category_id, name, description, price, preparation_time, 
		       is_available, is_vegetarian, is_vegan, is_gluten_free, spice_level, image_url
		FROM menu_items
	`
	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Close() error
	}
	var err error

	if categoryID != "" {
		query += " WHERE category_id = $1 ORDER BY name"
		rows, err = h.db.Query(query, categoryID)
	} else {
		query += " ORDER BY name"
		rows, err = h.db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var items []models.MenuItem
	for rows.Next() {
		var item models.MenuItem
		if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price,
			&item.PreparationTime, &item.IsAvailable, &item.IsVegetarian, &item.IsVegan,
			&item.IsGlutenFree, &item.SpiceLevel, &item.ImageURL); err != nil {
			continue
		}
		items = append(items, item)
	}

	// Sort items by ID for binary search price lookup optimization
	sort.Slice(items, func(i, j int) bool {
		return items[i].ID.String() < items[j].ID.String()
	})

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetMenuItem(c *gin.Context) {
	id := c.Param("id")

	var item models.MenuItem
	err := h.db.QueryRow(`
		SELECT id, category_id, name, description, price, preparation_time,
		       is_available, is_vegetarian, is_vegan, is_gluten_free, spice_level, image_url
		FROM menu_items WHERE id = $1
	`, id).Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price,
		&item.PreparationTime, &item.IsAvailable, &item.IsVegetarian, &item.IsVegan,
		&item.IsGlutenFree, &item.SpiceLevel, &item.ImageURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

type CreateMenuItemRequest struct {
	CategoryID      string  `json:"category_id" binding:"required"`
	Name            string  `json:"name" binding:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" binding:"required"`
	PreparationTime int     `json:"preparation_time"`
	IsVegetarian    bool    `json:"is_vegetarian"`
	IsVegan         bool    `json:"is_vegan"`
	IsGlutenFree    bool    `json:"is_gluten_free"`
	SpiceLevel      int     `json:"spice_level"`
	ImageURL        string  `json:"image_url"`
}

func (h *Handler) CreateMenuItem(c *gin.Context) {
	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var itemID uuid.UUID
	err := h.db.QueryRow(`
		INSERT INTO menu_items (category_id, name, description, price, preparation_time,
		                       is_vegetarian, is_vegan, is_gluten_free, spice_level, image_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`, req.CategoryID, req.Name, req.Description, req.Price, req.PreparationTime,
		req.IsVegetarian, req.IsVegan, req.IsGlutenFree, req.SpiceLevel, req.ImageURL).Scan(&itemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": itemID, "message": "Menu item created"})
}

func (h *Handler) UpdateMenuItem(c *gin.Context) {
	id := c.Param("id")

	var req CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE menu_items SET
			category_id = $1, name = $2, description = $3, price = $4, preparation_time = $5,
			is_vegetarian = $6, is_vegan = $7, is_gluten_free = $8, spice_level = $9, image_url = $10
		WHERE id = $11
	`, req.CategoryID, req.Name, req.Description, req.Price, req.PreparationTime,
		req.IsVegetarian, req.IsVegan, req.IsGlutenFree, req.SpiceLevel, req.ImageURL, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item updated"})
}

func (h *Handler) ToggleMenuItemAvailability(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE menu_items SET is_available = NOT is_available WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle availability"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Availability toggled"})
}

func (h *Handler) DeleteMenuItem(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE menu_items SET is_available = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete menu item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted"})
}

// Binary search for menu price lookup (used internally)
func binarySearchMenuPrice(items []models.MenuItem, targetID uuid.UUID) float64 {
	left, right := 0, len(items)-1

	for left <= right {
		mid := (left + right) / 2
		if items[mid].ID == targetID {
			return items[mid].Price
		}
		if items[mid].ID.String() < targetID.String() {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return 0
}
