package handlers

import (
	"net/http"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
)

// ========== REPORTS HANDLERS ==========

func (h *Handler) GetDashboardStats(c *gin.Context) {
	var stats models.DashboardStats

	// Total orders today
	h.db.QueryRow(`
		SELECT COUNT(*) FROM orders WHERE DATE(created_at) = CURRENT_DATE
	`).Scan(&stats.TotalOrders)

	// Active and total tables
	h.db.QueryRow(`
		SELECT COUNT(*) FROM restaurant_tables WHERE status = 'occupied'
	`).Scan(&stats.ActiveTables)

	h.db.QueryRow(`
		SELECT COUNT(*) FROM restaurant_tables
	`).Scan(&stats.TotalTables)

	// Daily revenue
	h.db.QueryRow(`
		SELECT COALESCE(SUM(total), 0) FROM orders 
		WHERE DATE(created_at) = CURRENT_DATE AND status = 'completed'
	`).Scan(&stats.DailyRevenue)

	// Pending and in-progress KOTs
	h.db.QueryRow(`
		SELECT COUNT(*) FROM kots WHERE status = 'pending'
	`).Scan(&stats.PendingKOTs)

	h.db.QueryRow(`
		SELECT COUNT(*) FROM kots WHERE status = 'in_progress'
	`).Scan(&stats.InProgressKOTs)

	c.JSON(http.StatusOK, stats)
}

func (h *Handler) GetSalesSummary(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", "")
	endDate := c.DefaultQuery("end_date", "")

	query := `
		SELECT 
			COUNT(*) as total_orders,
			COALESCE(SUM(total), 0) as total_revenue,
			COALESCE(AVG(total), 0) as avg_order_value,
			COUNT(DISTINCT table_id) as tables_served
		FROM orders
		WHERE status = 'completed'
	`
	args := []interface{}{}

	if startDate != "" {
		query += " AND created_at >= $1"
		args = append(args, startDate)
	}
	if endDate != "" {
		if len(args) > 0 {
			query += " AND created_at <= $2"
		} else {
			query += " AND created_at <= $1"
		}
		args = append(args, endDate)
	}

	var totalOrders int
	var totalRevenue, avgOrderValue float64
	var tablesServed int

	row := h.db.QueryRow(query, args...)
	row.Scan(&totalOrders, &totalRevenue, &avgOrderValue, &tablesServed)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":    totalOrders,
		"total_revenue":   totalRevenue,
		"avg_order_value": avgOrderValue,
		"tables_served":   tablesServed,
	})
}

func (h *Handler) GetPopularItems(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")

	rows, err := h.db.Query(`
		SELECT m.name, mc.name as category, COUNT(oi.id) as order_count, SUM(oi.total_price) as revenue
		FROM order_items oi
		JOIN menu_items m ON oi.menu_item_id = m.id
		LEFT JOIN menu_categories mc ON m.category_id = mc.id
		JOIN orders o ON oi.order_id = o.id
		WHERE o.status = 'completed'
		GROUP BY m.id, m.name, mc.name
		ORDER BY order_count DESC
		LIMIT $1
	`, limit)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var items []gin.H
	for rows.Next() {
		var name, category string
		var orderCount int
		var revenue float64
		if err := rows.Scan(&name, &category, &orderCount, &revenue); err != nil {
			continue
		}
		items = append(items, gin.H{
			"name":     name,
			"category": category,
			"orders":   orderCount,
			"revenue":  revenue,
		})
	}

	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetRevenueByDate(c *gin.Context) {
	days := c.DefaultQuery("days", "7")

	rows, err := h.db.Query(`
		SELECT DATE(created_at) as date, COUNT(*) as orders, COALESCE(SUM(total), 0) as revenue
		FROM orders
		WHERE status = 'completed' AND created_at >= CURRENT_DATE - INTERVAL '1 day' * $1::int
		GROUP BY DATE(created_at)
		ORDER BY date
	`, days)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var data []gin.H
	for rows.Next() {
		var date string
		var orders int
		var revenue float64
		if err := rows.Scan(&date, &orders, &revenue); err != nil {
			continue
		}
		data = append(data, gin.H{
			"date":    date,
			"orders":  orders,
			"revenue": revenue,
		})
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetOrdersByStatus(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT status, COUNT(*) as count
		FROM orders
		WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY status
		ORDER BY count DESC
	`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var total int
	var statusData []gin.H
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			continue
		}
		total += count
		statusData = append(statusData, gin.H{
			"status": status,
			"count":  count,
		})
	}

	// Calculate percentages
	for i := range statusData {
		if total > 0 {
			statusData[i]["percentage"] = float64(statusData[i]["count"].(int)) / float64(total) * 100
		}
	}

	c.JSON(http.StatusOK, statusData)
}
