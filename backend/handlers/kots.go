package handlers

import (
	"net/http"
	"sort"
	"time"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
)

// ========== KOT HANDLERS ==========

func (h *Handler) GetKOTs(c *gin.Context) {
	status := c.Query("status")

	query := `
		SELECT k.id, k.kot_number, k.order_id, o.order_number, COALESCE(t.table_number, ''),
		       k.status, k.priority, k.station, COALESCE(k.assigned_chef, ''), k.started_at, k.completed_at, k.created_at
		FROM kots k
		LEFT JOIN orders o ON k.order_id = o.id
		LEFT JOIN restaurant_tables t ON o.table_id = t.id
		WHERE 1=1
	`
	args := []interface{}{}

	if status != "" {
		query += " AND k.status = $1"
		args = append(args, status)
	}

	query += " ORDER BY k.priority DESC, k.created_at ASC"

	rows, err := h.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var kots []models.KOT
	for rows.Next() {
		var kot models.KOT
		if err := rows.Scan(&kot.ID, &kot.KOTNumber, &kot.OrderID, &kot.OrderNumber, &kot.TableNumber,
			&kot.Status, &kot.Priority, &kot.Station, &kot.AssignedChef, &kot.StartedAt,
			&kot.CompletedAt, &kot.CreatedAt); err != nil {
			continue
		}
		kots = append(kots, kot)
	}

	// Merge sort for stable priority ordering
	kots = mergeSortKOTs(kots)

	// Get items for each KOT
	for i := range kots {
		itemRows, err := h.db.Query(`
			SELECT id, menu_item_name, quantity, special_instructions, status
			FROM kot_items WHERE kot_id = $1
		`, kots[i].ID)
		if err == nil {
			for itemRows.Next() {
				var item models.KOTItem
				if err := itemRows.Scan(&item.ID, &item.MenuItemName, &item.Quantity,
					&item.SpecialInstructions, &item.Status); err != nil {
					continue
				}
				kots[i].Items = append(kots[i].Items, item)
			}
			itemRows.Close()
		}
	}

	c.JSON(http.StatusOK, kots)
}

// Merge sort for stable KOT ordering by priority and time
func mergeSortKOTs(kots []models.KOT) []models.KOT {
	if len(kots) <= 1 {
		return kots
	}

	mid := len(kots) / 2
	left := mergeSortKOTs(kots[:mid])
	right := mergeSortKOTs(kots[mid:])

	return mergeKOTs(left, right)
}

func mergeKOTs(left, right []models.KOT) []models.KOT {
	result := make([]models.KOT, 0, len(left)+len(right))
	i, j := 0, 0

	for i < len(left) && j < len(right) {
		// Higher priority first, then earlier creation time
		if left[i].Priority > right[j].Priority ||
			(left[i].Priority == right[j].Priority && left[i].CreatedAt.Before(right[j].CreatedAt)) {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}

	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func (h *Handler) GetKOT(c *gin.Context) {
	id := c.Param("id")

	var kot models.KOT
	err := h.db.QueryRow(`
		SELECT k.id, k.kot_number, k.order_id, o.order_number, COALESCE(t.table_number, ''),
		       k.status, k.priority, k.station, COALESCE(k.assigned_chef, ''), k.started_at, k.completed_at, k.created_at
		FROM kots k
		LEFT JOIN orders o ON k.order_id = o.id
		LEFT JOIN restaurant_tables t ON o.table_id = t.id
		WHERE k.id = $1
	`, id).Scan(&kot.ID, &kot.KOTNumber, &kot.OrderID, &kot.OrderNumber, &kot.TableNumber,
		&kot.Status, &kot.Priority, &kot.Station, &kot.AssignedChef, &kot.StartedAt,
		&kot.CompletedAt, &kot.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "KOT not found"})
		return
	}

	// Get items
	itemRows, _ := h.db.Query(`
		SELECT id, menu_item_name, quantity, special_instructions, status
		FROM kot_items WHERE kot_id = $1
	`, id)
	defer itemRows.Close()

	for itemRows.Next() {
		var item models.KOTItem
		if err := itemRows.Scan(&item.ID, &item.MenuItemName, &item.Quantity,
			&item.SpecialInstructions, &item.Status); err != nil {
			continue
		}
		kot.Items = append(kot.Items, item)
	}

	c.JSON(http.StatusOK, kot)
}

func (h *Handler) UpdateKOTStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var startedAt, completedAt interface{}
	if req.Status == "in_progress" {
		startedAt = time.Now()
	} else if req.Status == "completed" {
		completedAt = time.Now()
	}

	_, err := h.db.Exec(`
		UPDATE kots SET status = $1, started_at = COALESCE($2, started_at), completed_at = COALESCE($3, completed_at)
		WHERE id = $4
	`, req.Status, startedAt, completedAt, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update KOT status"})
		return
	}

	// If completed, update order status
	if req.Status == "completed" {
		h.db.Exec(`
			UPDATE orders SET status = 'ready'
			WHERE id = (SELECT order_id FROM kots WHERE id = $1)
		`, id)
	}

	c.JSON(http.StatusOK, gin.H{"message": "KOT status updated"})
}

func (h *Handler) AssignKOTChef(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		AssignedChef string `json:"assigned_chef" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec("UPDATE kots SET assigned_chef = $1 WHERE id = $2", req.AssignedChef, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign chef"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chef assigned"})
}

func (h *Handler) UpdateKOTItemStatus(c *gin.Context) {
	itemID := c.Param("itemId")

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.db.Exec("UPDATE kot_items SET status = $1 WHERE id = $2", req.Status, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item status updated"})
}

// ========== PAYMENT HANDLERS ==========

func (h *Handler) GetPayments(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT p.id, p.order_id, p.payment_method_id, pm.name, p.amount, p.tip, p.status,
		       p.transaction_reference, p.processed_at, p.created_at
		FROM payments p
		LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
		ORDER BY p.created_at DESC
		LIMIT 100
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		if err := rows.Scan(&p.ID, &p.OrderID, &p.PaymentMethodID, &p.PaymentMethodName,
			&p.Amount, &p.Tip, &p.Status, &p.TransactionReference, &p.ProcessedAt, &p.CreatedAt); err != nil {
			continue
		}
		payments = append(payments, p)
	}

	c.JSON(http.StatusOK, payments)
}

func (h *Handler) GetPayment(c *gin.Context) {
	id := c.Param("id")

	var p models.Payment
	err := h.db.QueryRow(`
		SELECT p.id, p.order_id, p.payment_method_id, pm.name, p.amount, p.tip, p.status,
		       p.transaction_reference, p.processed_at, p.created_at
		FROM payments p
		LEFT JOIN payment_methods pm ON p.payment_method_id = pm.id
		WHERE p.id = $1
	`, id).Scan(&p.ID, &p.OrderID, &p.PaymentMethodID, &p.PaymentMethodName,
		&p.Amount, &p.Tip, &p.Status, &p.TransactionReference, &p.ProcessedAt, &p.CreatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

type CreatePaymentRequest struct {
	OrderID         string  `json:"order_id" binding:"required"`
	PaymentMethodID string  `json:"payment_method_id" binding:"required"`
	Amount          float64 `json:"amount" binding:"required"`
	Tip             float64 `json:"tip"`
}

func (h *Handler) CreatePayment(c *gin.Context) {
	var req CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")
	now := time.Now()

	// Generate transaction reference
	txRef, _ := generateSecureToken(16)

	var paymentID string
	err := h.db.QueryRow(`
		INSERT INTO payments (order_id, payment_method_id, amount, tip, status, transaction_reference, processed_by, processed_at)
		VALUES ($1, $2, $3, $4, 'completed', $5, $6, $7)
		RETURNING id
	`, req.OrderID, req.PaymentMethodID, req.Amount, req.Tip, txRef, userID, now).Scan(&paymentID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	// Update order status to completed
	h.db.Exec("UPDATE orders SET status = 'completed' WHERE id = $1", req.OrderID)

	c.JSON(http.StatusCreated, gin.H{
		"id":                    paymentID,
		"transaction_reference": txRef,
		"message":               "Payment processed successfully",
	})
}

func (h *Handler) RefundPayment(c *gin.Context) {
	id := c.Param("id")

	_, err := h.db.Exec("UPDATE payments SET status = 'refunded' WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refund payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment refunded"})
}

func (h *Handler) GetPaymentMethods(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, name, is_active FROM payment_methods WHERE is_active = true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var methods []models.PaymentMethod
	for rows.Next() {
		var m models.PaymentMethod
		if err := rows.Scan(&m.ID, &m.Name, &m.IsActive); err != nil {
			continue
		}
		methods = append(methods, m)
	}

	// Sort for consistent ordering
	sort.Slice(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})

	c.JSON(http.StatusOK, methods)
}
