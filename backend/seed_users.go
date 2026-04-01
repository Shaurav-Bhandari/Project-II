//go:build ignore

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	godotenv.Load()

	host := envOr("DB_HOST", "localhost")
	port := envOr("DB_PORT", "5432")
	user := envOr("DB_USER", "postgres")
	pass := envOr("DB_PASSWORD", "postgres")
	name := envOr("DB_NAME", "restaurant_db")
	ssl := envOr("DB_SSLMODE", "disable")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, name, ssl)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB open error:", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("DB ping error:", err)
	}
	fmt.Println("✅ Connected to database")

	password := "staff123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Bcrypt error:", err)
	}
	hashStr := string(hash)

	// Verify hash before inserting
	if err := bcrypt.CompareHashAndPassword([]byte(hashStr), []byte(password)); err != nil {
		log.Fatal("Hash verification failed:", err)
	}
	fmt.Println("✅ Hash verified for password:", password)

	// Get role IDs
	var staffRoleID, kitchenRoleID, managerRoleID string
	db.QueryRow("SELECT id FROM roles WHERE name = 'staff'").Scan(&staffRoleID)
	db.QueryRow("SELECT id FROM roles WHERE name = 'kitchen'").Scan(&kitchenRoleID)
	db.QueryRow("SELECT id FROM roles WHERE name = 'manager'").Scan(&managerRoleID)

	// Delete existing seeded users (preserve admin)
	db.Exec("DELETE FROM users WHERE email IN ('waiter1@restaurant.com','waiter2@restaurant.com','waiter3@restaurant.com','kitchen1@restaurant.com','kitchen2@restaurant.com','manager@restaurant.com')")
	fmt.Println("✅ Cleared existing seeded users")

	users := []struct {
		email, first, last, roleID string
	}{
		{"waiter1@restaurant.com", "Ram", "Thapa", staffRoleID},
		{"waiter2@restaurant.com", "Sita", "Gurung", staffRoleID},
		{"waiter3@restaurant.com", "Hari", "Tamang", staffRoleID},
		{"kitchen1@restaurant.com", "Bishnu", "Sherpa", kitchenRoleID},
		{"kitchen2@restaurant.com", "Maya", "Rai", kitchenRoleID},
		{"manager@restaurant.com", "Prakash", "Shrestha", managerRoleID},
	}

	for _, u := range users {
		_, err := db.Exec(
			"INSERT INTO users (email, password_hash, first_name, last_name, role_id, is_active) VALUES ($1, $2, $3, $4, $5, true)",
			u.email, hashStr, u.first, u.last, u.roleID,
		)
		if err != nil {
			log.Printf("⚠️  Failed to insert %s: %v", u.email, err)
		} else {
			fmt.Printf("✅ Created user: %s (%s %s)\n", u.email, u.first, u.last)
		}
	}

	// Verify by reading back and comparing
	var dbHash string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = 'waiter1@restaurant.com'").Scan(&dbHash)
	if err != nil {
		log.Fatal("Failed to read back hash:", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password)); err != nil {
		log.Fatal("❌ ROUND-TRIP VERIFICATION FAILED:", err)
	}
	fmt.Println("✅ Round-trip verification passed — login should work!")
	fmt.Println("\nCredentials: all seeded users use password 'staff123'")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
