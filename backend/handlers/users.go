package handlers

import (
	"net/http"

	"restaurant-management/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ========== USER HANDLERS ==========

func (h *Handler) GetUsers(c *gin.Context) {
	rows, err := h.db.Query(`
		SELECT u.id, u.email, u.first_name, u.last_name, u.is_active, u.last_login, u.created_at, r.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var users []gin.H
	for rows.Next() {
		var user models.User
		var roleName string
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.IsActive, &user.LastLogin, &user.CreatedAt, &roleName); err != nil {
			continue
		}
		users = append(users, gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"is_active":  user.IsActive,
			"last_login": user.LastLogin,
			"created_at": user.CreatedAt,
			"role":       roleName,
		})
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	var roleName string
	err := h.db.QueryRow(`
		SELECT u.id, u.email, u.first_name, u.last_name, u.is_active, u.last_login, r.name
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = $1
	`, id).Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.IsActive, &user.LastLogin, &roleName)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"is_active":  user.IsActive,
		"last_login": user.LastLogin,
		"role":       roleName,
	})
}

type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	var roleID uuid.UUID
	err = h.db.QueryRow("SELECT id FROM roles WHERE name = $1", req.Role).Scan(&roleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	var userID uuid.UUID
	err = h.db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, req.Email, string(hashedPassword), req.FirstName, req.LastName, roleID).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": userID, "message": "User created successfully"})
}

type UpdateUserRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	IsActive  *bool  `json:"is_active"`
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var roleID uuid.UUID
	if req.Role != "" {
		err := h.db.QueryRow("SELECT id FROM roles WHERE name = $1", req.Role).Scan(&roleID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
	}

	_, err := h.db.Exec(`
		UPDATE users SET
			email = COALESCE(NULLIF($1, ''), email),
			first_name = COALESCE(NULLIF($2, ''), first_name),
			last_name = COALESCE(NULLIF($3, ''), last_name),
			role_id = CASE WHEN $4 != '00000000-0000-0000-0000-000000000000' THEN $4::uuid ELSE role_id END,
			is_active = COALESCE($5, is_active)
		WHERE id = $6
	`, req.Email, req.FirstName, req.LastName, roleID.String(), req.IsActive, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	// Soft delete by deactivating
	_, err := h.db.Exec("UPDATE users SET is_active = false WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deactivated successfully"})
}

func (h *Handler) GetRoles(c *gin.Context) {
	rows, err := h.db.Query("SELECT id, name, description FROM roles ORDER BY name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			continue
		}
		roles = append(roles, role)
	}

	c.JSON(http.StatusOK, roles)
}
