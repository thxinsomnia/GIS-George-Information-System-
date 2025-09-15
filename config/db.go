package config

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
        log.Fatal("Database URL is Not Set in Environment Variables")
    }
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to Connect to Database: %v", err)
	}

    log.Println("Database Connection Successful.")
	DB = db
}

var SupabaseClient *supabase.Client

// Fungsi untuk inisialisasi client (panggil sekali saat aplikasi start)
func InitSupabase() {
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		panic("Cannot initialize Supabase client")
	}
	SupabaseClient = client
}
