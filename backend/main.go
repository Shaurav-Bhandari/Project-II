package main

import (
	"log"
	"os"

	"restaurant-management/config"
	"restaurant-management/handlers"
	"restaurant-management/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	
	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize admin user if not exists
	if err := config.InitializeAdmin(db); err != nil {
		log.Printf("Warning: Failed to initialize admin user: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Public routes
	public := r.Group("/api")
	{
		public.POST("/auth/login", h.Login)
		public.POST("/auth/register", h.Register)
	}

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		// Auth
		api.POST("/auth/logout", h.Logout)
		api.GET("/auth/me", h.GetCurrentUser)

		// Users (admin only)
		users := api.Group("/users")
		users.Use(middleware.RoleMiddleware("admin"))
		{
			users.GET("", h.GetUsers)
			users.GET("/:id", h.GetUser)
			users.POST("", h.CreateUser)
			users.PUT("/:id", h.UpdateUser)
			users.DELETE("/:id", h.DeleteUser)
		}

		// Roles
		api.GET("/roles", h.GetRoles)

		// Tables
		tables := api.Group("/tables")
		{
			tables.GET("", h.GetTables)
			tables.GET("/:id", h.GetTable)
			tables.POST("", h.CreateTable)
			tables.PUT("/:id", h.UpdateTable)
			tables.PATCH("/:id/status", h.UpdateTableStatus)
			tables.DELETE("/:id", h.DeleteTable)
			tables.POST("/allocate", h.AllocateTables)
		}

		// Menu
		menu := api.Group("/menu")
		{
			menu.GET("/categories", h.GetCategories)
			menu.POST("/categories", h.CreateCategory)
			menu.PUT("/categories/:id", h.UpdateCategory)
			menu.DELETE("/categories/:id", h.DeleteCategory)

			menu.GET("/items", h.GetMenuItems)
			menu.GET("/items/:id", h.GetMenuItem)
			menu.POST("/items", h.CreateMenuItem)
			menu.PUT("/items/:id", h.UpdateMenuItem)
			menu.PATCH("/items/:id/availability", h.ToggleMenuItemAvailability)
			menu.DELETE("/items/:id", h.DeleteMenuItem)
		}

		// Orders
		orders := api.Group("/orders")
		{
			orders.GET("", h.GetOrders)
			orders.GET("/:id", h.GetOrder)
			orders.POST("", h.CreateOrder)
			orders.PUT("/:id", h.UpdateOrder)
			orders.PATCH("/:id/status", h.UpdateOrderStatus)
			orders.DELETE("/:id", h.CancelOrder)
			orders.GET("/:id/receipt", h.GenerateReceipt)

			orders.POST("/:id/items", h.AddOrderItem)
			orders.PUT("/:id/items/:itemId", h.UpdateOrderItem)
			orders.DELETE("/:id/items/:itemId", h.RemoveOrderItem)
		}

		// KOTs
		kots := api.Group("/kots")
		{
			kots.GET("", h.GetKOTs)
			kots.GET("/:id", h.GetKOT)
			kots.PATCH("/:id/status", h.UpdateKOTStatus)
			kots.PATCH("/:id/assign", h.AssignKOTChef)
			kots.PATCH("/:id/items/:itemId/status", h.UpdateKOTItemStatus)
		}

		// Payments
		payments := api.Group("/payments")
		{
			payments.GET("", h.GetPayments)
			payments.GET("/:id", h.GetPayment)
			payments.POST("", h.CreatePayment)
			payments.POST("/:id/refund", h.RefundPayment)
		}
		api.GET("/payment-methods", h.GetPaymentMethods)

		// Reports
		reports := api.Group("/reports")
		reports.Use(middleware.RoleMiddleware("admin", "manager"))
		{
			reports.GET("/dashboard", h.GetDashboardStats)
			reports.GET("/sales", h.GetSalesSummary)
			reports.GET("/popular-items", h.GetPopularItems)
			reports.GET("/revenue", h.GetRevenueByDate)
			reports.GET("/orders-by-status", h.GetOrdersByStatus)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
