package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Cfg struct {
	DBConnectionString string
	SupabaseURL        string
	SupabaseKey        string
}

func Get() Cfg {
	return Cfg{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseKey:        os.Getenv("SUPABASE_KEY"),
	}
}
