package config

import (
	"log"
	"os"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	godotenv.Load() // Load .env file

	// 1. Get each variable from the .env file
	dbUser := os.Getenv("user")
	dbPass := os.Getenv("password")
	dbHost := os.Getenv("host")
	dbPort := os.Getenv("port")
	dbName := os.Getenv("dbname")

	// 2. Construct the DSN string for the Connection Pooler
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require prefer_simple_protocol=true",
		dbHost, dbUser, dbPass, dbName, dbPort)

	// 3. Connect using the constructed DSN
	db, err := gorm.Open(postgres.New(postgres.Config{
    DSN: dsn,
    PreferSimpleProtocol: true, // disables prepared statement cache
}), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
	
	DB = db
	log.Println("Database connection successful via Connection Pooler.")
}

var SupabaseClient *supabase.Client

// Fungsi untuk inisialisasi client (panggil sekali saat aplikasi start)
func InitSupabase() {
	supabaseURL := os.Getenv("DATABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SECRET_KEY")
	client, err := supabase.NewClient(supabaseURL, supabaseKey, nil)
	if err != nil {
		panic("Cannot initialize Supabase client")
	}
	SupabaseClient = client
}

