package handlers

import (
	"net/http"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ========== ORDER HANDLERS ==========

func (h *Handler) GetOrders(c *gin.Context) {
	status := c.Query("status")
	tableID := c.Query("table_id")

	query := `
		SELECT o.id, o.order_number, o.table_id, COALESCE(t.table_number, ''), o.customer_name,
		       o.status, o.order_type, o.notes, o.subtotal, o.tax, o.discount, o.total, o.created_at
		FROM orders o
		LEFT JOIN restaurant_tables t ON o.table_id = t.id
		WHERE 1=1
	`
	args := []interface{}{}
	argCount := 0

	if status != "" {
		argCount++
		query += " AND o.status = $" + string(rune('0'+argCount))
		args = append(args, status)
	}
	if tableID != "" {
		argCount++
		query += " AND o.table_id = $" + string(rune('0'+argCount))
		args = append(args, tableID)
	}

	query += " ORDER BY o.created_at DESC LIMIT 100"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.OrderNumber, &order.TableID, &order.TableNumber,
			&order.CustomerName, &order.Status, &order.OrderType, &order.Notes, &order.Subtotal,
			&order.Tax, &order.Discount, &order.Total, &order.CreatedAt); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	err := h.db.QueryRow(`
		SELECT o.id, o.order_number, o.table_id, COALESCE(t.table_number, ''), o.customer_name,
		       o.status, o.order_type, o.notes, o.subtotal, o.tax, o.discount, o.total, o.created_at
		FROM orders o
		LEFT JOIN restaurant_tables t ON o.table_id = t.id
		WHERE o.id = $1
	`, id).Scan(&order.ID, &order.OrderNumber, &order.TableID, &order.TableNumber,
		&order.CustomerName, &order.Status, &order.OrderType, &order.Notes, &order.Subtotal,
		&order.Tax, &order.Discount, &order.Total, &order.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Get order items
	itemRows, err := h.db.Query(`
		SELECT oi.id, oi.menu_item_id, m.name, oi.quantity, oi.unit_price, oi.total_price,
		       oi.special_instructions, oi.status
		FROM order_items oi
		LEFT JOIN menu_items m ON oi.menu_item_id = m.id
		WHERE oi.order_id = $1
	`, id)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item models.OrderItem
			if err := itemRows.Scan(&item.ID, &item.MenuItemID, &item.MenuItemName, &item.Quantity,
				&item.UnitPrice, &item.TotalPrice, &item.SpecialInstructions, &item.Status); err != nil {
				continue
			}
			order.Items = append(order.Items, item)
		}
	}

	c.JSON(http.StatusOK, order)
}

type CreateOrderRequest struct {
	TableID      string           `json:"table_id"`
	CustomerName string           `json:"customer_name"`
	OrderType    string           `json:"order_type"`
	Notes        string           `json:"notes"`
	Items        []OrderItemInput `json:"items" binding:"required,min=1"`
}

type OrderItemInput struct {
	MenuItemID          string `json:"menu_item_id" binding:"required"`
	Quantity            int    `json:"quantity" binding:"required,min=1"`
	SpecialInstructions string `json:"special_instructions"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	orderType := req.OrderType
	if orderType == "" {
		orderType = "dine-in"
	}

	// Start transaction
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Create order
	var orderID uuid.UUID
	var tableID interface{} = nil
	if req.TableID != "" {
		tableID = req.TableID
	}

	err = tx.QueryRow(`
		INSERT INTO orders (table_id, customer_name, order_type, notes, created_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, tableID, req.CustomerName, orderType, req.Notes, userID).Scan(&orderID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Add order items
	for _, item := range req.Items {
		var price float64
		err := h.db.QueryRow("SELECT price FROM menu_items WHERE id = $1", item.MenuItemID).Scan(&price)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item"})
			return
		}

		totalPrice := price * float64(item.Quantity)
		_, err = tx.Exec(`
			INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price, total_price, special_instructions)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, orderID, item.MenuItemID, item.Quantity, price, totalPrice, item.SpecialInstructions)

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add order items"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": orderID, "message": "Order created successfully"})
}

func (h *Handler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CustomerName string `json:"customer_name"`
		Notes        string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE orders SET customer_name = $1, notes = $2 WHERE id = $3
	`, req.CustomerName, req.Notes, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order updated"})
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// The database trigger will auto-generate KOT when status changes to 'confirmed'
	_, err := h.db.Exec("UPDATE orders SET status = $1 WHERE id = $2", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func (h *Handler) CancelOrder(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE orders SET status = 'cancelled' WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order cancelled"})
}

func (h *Handler) AddOrderItem(c *gin.Context) {
	orderID := c.Param("id")

	var req OrderItemInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var price float64
	err := h.db.QueryRow("SELECT price FROM menu_items WHERE id = $1", req.MenuItemID).Scan(&price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item"})
		return
	}

	totalPrice := price * float64(req.Quantity)
	var itemID uuid.UUID
	err = h.db.QueryRow(`
		INSERT INTO order_items (order_id, menu_item_id, quantity, unit_price, total_price, special_instructions)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, orderID, req.MenuItemID, req.Quantity, price, totalPrice, req.SpecialInstructions).Scan(&itemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": itemID, "message": "Item added"})
}

func (h *Handler) UpdateOrderItem(c *gin.Context) {
	itemID := c.Param("itemId")

	var req struct {
		Quantity            int    `json:"quantity"`
		SpecialInstructions string `json:"special_instructions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec(`
		UPDATE order_items SET quantity = $1, total_price = unit_price * $1, special_instructions = $2
		WHERE id = $3
	`, req.Quantity, req.SpecialInstructions, itemID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item updated"})
}

func (h *Handler) RemoveOrderItem(c *gin.Context) {
	itemID := c.Param("itemId")

	_, err := h.db.Exec("DELETE FROM order_items WHERE id = $1", itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed"})
}

func (h *Handler) GenerateReceipt(c *gin.Context) {
	id := c.Param("id")

	var order models.Order
	err := h.db.QueryRow(`
		SELECT o.id, o.order_number, COALESCE(t.table_number, 'Takeaway'), o.customer_name,
		       o.subtotal, o.tax, o.discount, o.total, o.created_at
		FROM orders o
		LEFT JOIN restaurant_tables t ON o.table_id = t.id
		WHERE o.id = $1
	`, id).Scan(&order.ID, &order.OrderNumber, &order.TableNumber, &order.CustomerName,
		&order.Subtotal, &order.Tax, &order.Discount, &order.Total, &order.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Get items
	itemRows, _ := h.db.Query(`
		SELECT m.name, oi.quantity, oi.unit_price, oi.total_price
		FROM order_items oi
		LEFT JOIN menu_items m ON oi.menu_item_id = m.id
		WHERE oi.order_id = $1
	`, id)
	defer itemRows.Close()

	var items []gin.H
	for itemRows.Next() {
		var name string
		var qty int
		var unitPrice, totalPrice float64
		itemRows.Scan(&name, &qty, &unitPrice, &totalPrice)
		items = append(items, gin.H{
			"name":        name,
			"quantity":    qty,
			"unit_price":  unitPrice,
			"total_price": totalPrice,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"order_number":  order.OrderNumber,
		"table":         order.TableNumber,
		"customer_name": order.CustomerName,
		"items":         items,
		"subtotal":      order.Subtotal,
		"tax":           order.Tax,
		"discount":      order.Discount,
		"total":         order.Total,
		"date":          order.CreatedAt,
	})
}
