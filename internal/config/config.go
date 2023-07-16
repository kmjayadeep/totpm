package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Cfg struct {
	DBConnectionString string
	SupabaseURL        string
	SupabaseKey        string
	JWTSecret          string
}

func Get() Cfg {
	return Cfg{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseKey:        os.Getenv("SUPABASE_KEY"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
	}
}
