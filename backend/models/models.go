package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a system user
type User struct {
	ID           uuid.UUID  `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	RoleID       uuid.UUID  `json:"role_id"`
	Role         *Role      `json:"role,omitempty"`
	IsActive     bool       `json:"is_active"`
	LastLogin    *time.Time `json:"last_login"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Role represents a user role
type Role struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Permissions string    `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
}

// RestaurantTable represents a restaurant table
type RestaurantTable struct {
	ID          uuid.UUID `json:"id"`
	TableNumber string    `json:"table_number"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"`
	Location    string    `json:"location"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MenuCategory represents a menu category
type MenuCategory struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// MenuItem represents a menu item
type MenuItem struct {
	ID              uuid.UUID `json:"id"`
	CategoryID      uuid.UUID `json:"category_id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Price           float64   `json:"price"`
	PreparationTime int       `json:"preparation_time"`
	IsAvailable     bool      `json:"is_available"`
	IsVegetarian    bool      `json:"is_vegetarian"`
	IsVegan         bool      `json:"is_vegan"`
	IsGlutenFree    bool      `json:"is_gluten_free"`
	SpiceLevel      int       `json:"spice_level"`
	ImageURL        string    `json:"image_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Order represents a customer order
type Order struct {
	ID           uuid.UUID   `json:"id"`
	OrderNumber  int         `json:"order_number"`
	TableID      *uuid.UUID  `json:"table_id"`
	TableNumber  string      `json:"table_number,omitempty"`
	CustomerName string      `json:"customer_name"`
	Status       string      `json:"status"`
	OrderType    string      `json:"order_type"`
	Notes        string      `json:"notes"`
	Subtotal     float64     `json:"subtotal"`
	Tax          float64     `json:"tax"`
	Discount     float64     `json:"discount"`
	Total        float64     `json:"total"`
	CreatedBy    uuid.UUID   `json:"created_by"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Items        []OrderItem `json:"items,omitempty"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID                  uuid.UUID `json:"id"`
	OrderID             uuid.UUID `json:"order_id"`
	MenuItemID          uuid.UUID `json:"menu_item_id"`
	MenuItemName        string    `json:"menu_item_name,omitempty"`
	Quantity            int       `json:"quantity"`
	UnitPrice           float64   `json:"unit_price"`
	TotalPrice          float64   `json:"total_price"`
	SpecialInstructions string    `json:"special_instructions"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// KOT represents a Kitchen Order Ticket
type KOT struct {
	ID           uuid.UUID  `json:"id"`
	KOTNumber    int        `json:"kot_number"`
	OrderID      uuid.UUID  `json:"order_id"`
	OrderNumber  int        `json:"order_number,omitempty"`
	TableNumber  string     `json:"table_number,omitempty"`
	Status       string     `json:"status"`
	Priority     int        `json:"priority"`
	Station      string     `json:"station"`
	AssignedChef string     `json:"assigned_chef"`
	StartedAt    *time.Time `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Items        []KOTItem  `json:"items,omitempty"`
}

// KOTItem represents an item in a KOT
type KOTItem struct {
	ID                  uuid.UUID `json:"id"`
	KOTID               uuid.UUID `json:"kot_id"`
	OrderItemID         uuid.UUID `json:"order_item_id"`
	MenuItemName        string    `json:"menu_item_name"`
	Quantity            int       `json:"quantity"`
	SpecialInstructions string    `json:"special_instructions"`
	Status              string    `json:"status"`
	CreatedAt           time.Time `json:"created_at"`
}

// Payment represents a payment record
type Payment struct {
	ID                   uuid.UUID  `json:"id"`
	OrderID              uuid.UUID  `json:"order_id"`
	PaymentMethodID      uuid.UUID  `json:"payment_method_id"`
	PaymentMethodName    string     `json:"payment_method_name,omitempty"`
	Amount               float64    `json:"amount"`
	Tip                  float64    `json:"tip"`
	Status               string     `json:"status"`
	TransactionReference string     `json:"transaction_reference"`
	ProcessedBy          uuid.UUID  `json:"processed_by"`
	ProcessedAt          *time.Time `json:"processed_at"`
	CreatedAt            time.Time  `json:"created_at"`
}

// PaymentMethod represents a payment method
type PaymentMethod struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Dashboard stats
type DashboardStats struct {
	TotalOrders    int     `json:"total_orders"`
	ActiveTables   int     `json:"active_tables"`
	TotalTables    int     `json:"total_tables"`
	DailyRevenue   float64 `json:"daily_revenue"`
	PendingKOTs    int     `json:"pending_kots"`
	InProgressKOTs int     `json:"in_progress_kots"`
}
