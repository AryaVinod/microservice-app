package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Successfully connected to database!")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go Microservice! 🚀")
}

func dbTestHandler(w http.ResponseWriter, r *http.Request) {
	// Create a simple table if it doesn't exist
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS visits (
		id SERIAL PRIMARY KEY,
		visited_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating table: %v", err), 500)
		return
	}

	// Insert a new visit
	_, err = db.Exec(`INSERT INTO visits DEFAULT VALUES`)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting visit: %v", err), 500)
		return
	}

	// Count total visits
	var count int
	err = db.QueryRow(`SELECT COUNT(*) FROM visits`).Scan(&count)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error counting visits: %v", err), 500)
		return
	}

	fmt.Fprintf(w, "✅ Database connected! Total visits: %d", count)
}

func main() {
	connectDB()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/db-test", dbTestHandler)

	port := "8080"
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}