package config

import (
	"database/sql"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

// InitializeAdmin creates a default admin user if one doesn't exist
func InitializeAdmin(db *sql.DB) error {
	// Check if any admin user exists
	var count int
	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE r.name = 'admin'
	`).Scan(&count)

	if err != nil {
		return err
	}

	// If admin exists, skip initialization
	if count > 0 {
		log.Println("Admin user already exists, skipping initialization")
		return nil
	}

	// Get admin credentials from environment or use defaults
	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		adminEmail = "admin@restaurant.com"
	}

	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		adminPassword = "admin123" // Default password
		log.Println("⚠️  WARNING: Using default admin password. Set ADMIN_PASSWORD environment variable for production!")
	}

	adminFirstName := os.Getenv("ADMIN_FIRST_NAME")
	if adminFirstName == "" {
		adminFirstName = "System"
	}

	adminLastName := os.Getenv("ADMIN_LAST_NAME")
	if adminLastName == "" {
		adminLastName = "Admin"
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Get admin role ID
	var roleID string
	err = db.QueryRow("SELECT id FROM roles WHERE name = 'admin'").Scan(&roleID)
	if err != nil {
		return err
	}

	// Create admin user
	var userID string
	err = db.QueryRow(`
		INSERT INTO users (email, password_hash, first_name, last_name, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, true)
		RETURNING id
	`, adminEmail, string(hashedPassword), adminFirstName, adminLastName, roleID).Scan(&userID)

	if err != nil {
		return err
	}

	log.Printf("✅ Admin user created successfully!")
	log.Printf("   Email: %s", adminEmail)
	log.Printf("   Password: %s", adminPassword)
	log.Printf("   User ID: %s", userID)
	log.Println("   ⚠️  Please change the password after first login!")

	return nil
}
