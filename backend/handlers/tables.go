package handlers

import (
	"net/http"
	"sort"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ========== TABLE HANDLERS ==========

func (h *Handler) GetTables(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT id, table_number, capacity, status, location, created_at
		FROM restaurant_tables
		ORDER BY table_number
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var tables []models.RestaurantTable
	for rows.Next() {
		var table models.RestaurantTable
		if err := rows.Scan(&table.ID, &table.TableNumber, &table.Capacity, &table.Status, &table.Location, &table.CreatedAt); err != nil {
			continue
		}
		tables = append(tables, table)
	}

	c.JSON(http.StatusOK, tables)
}

func (h *Handler) GetTable(c *gin.Context) {
	id := c.Param("id")

	var table models.RestaurantTable
	err := h.db.QueryRow(`
		SELECT id, table_number, capacity, status, location, created_at
		FROM restaurant_tables WHERE id = $1
	`, id).Scan(&table.ID, &table.TableNumber, &table.Capacity, &table.Status, &table.Location, &table.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, table)
}

type CreateTableRequest struct {
	TableNumber string `json:"table_number" binding:"required"`
	Capacity    int    `json:"capacity" binding:"required,min=1"`
	Location    string `json:"location"`
}

func (h *Handler) CreateTable(c *gin.Context) {
	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tableID uuid.UUID
	err := h.db.QueryRow(`
		INSERT INTO restaurant_tables (table_number, capacity, location)
		VALUES ($1, $2, $3)
		RETURNING id
	`, req.TableNumber, req.Capacity, req.Location).Scan(&tableID)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Table number already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": tableID, "message": "Table created successfully"})
}

func (h *Handler) UpdateTable(c *gin.Context) {
	id := c.Param("id")

	var req CreateTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE restaurant_tables SET
			table_number = $1, capacity = $2, location = $3
		WHERE id = $4
	`, req.TableNumber, req.Capacity, req.Location, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully"})
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateTableStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec("UPDATE restaurant_tables SET status = $1 WHERE id = $2", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status updated successfully"})
}

func (h *Handler) DeleteTable(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("DELETE FROM restaurant_tables WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
}

// Backtracking algorithm for table allocation
type AllocateRequest struct {
	GroupSize int `json:"group_size" binding:"required,min=1"`
}

func (h *Handler) AllocateTables(c *gin.Context) {
	var req AllocateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get available tables
	rows, err := h.db.Query(`
		SELECT id, table_number, capacity
		FROM restaurant_tables
		WHERE status = 'available'
		ORDER BY capacity DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var tables []models.RestaurantTable
	for rows.Next() {
		var table models.RestaurantTable
		if err := rows.Scan(&table.ID, &table.TableNumber, &table.Capacity); err != nil {
			continue
		}
		tables = append(tables, table)
	}

	// Sort by capacity descending for greedy approach
	sort.Slice(tables, func(i, j int) bool {
		return tables[i].Capacity > tables[j].Capacity
	})

	// Backtracking to find optimal combination
	result := findTableCombination(tables, req.GroupSize, []models.RestaurantTable{}, 0)

	if result == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "No available table combination found",
			"message": "Not enough table capacity for the group size",
		})
		return
	}

	var tableNumbers []string
	var totalCapacity int
	for _, t := range result {
		tableNumbers = append(tableNumbers, t.TableNumber)
		totalCapacity += t.Capacity
	}

	c.JSON(http.StatusOK, gin.H{
		"tables":         tableNumbers,
		"total_capacity": totalCapacity,
		"group_size":     req.GroupSize,
	})
}

// Backtracking helper function
func findTableCombination(tables []models.RestaurantTable, target int, current []models.RestaurantTable, start int) []models.RestaurantTable {
	currentCapacity := 0
	for _, t := range current {
		currentCapacity += t.Capacity
	}

	// Found a valid combination
	if currentCapacity >= target {
		result := make([]models.RestaurantTable, len(current))
		copy(result, current)
		return result
	}

	// No more tables to try
	if start >= len(tables) {
		return nil
	}

	// Try including current table
	withCurrent := findTableCombination(tables, target, append(current, tables[start]), start+1)
	if withCurrent != nil {
		return withCurrent
	}

	// Try without current table (backtrack)
	return findTableCombination(tables, target, current, start+1)
}
