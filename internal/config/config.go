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
	AppKey             string
	EnableDbMigration  bool
	S3Key              string
	S3Secret           string
	S3Bucket           string
	S3Endpoint         string
	S3Region           string
}

func Get() Cfg {
	return Cfg{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseKey:        os.Getenv("SUPABASE_KEY"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		AppKey:             os.Getenv("APP_KEY"),
		EnableDbMigration:  os.Getenv("ENABLE_DB_MIGRATION") == "true",
		S3Key:              os.Getenv("AWS_ACCESS_KEY_ID"),
		S3Secret:           os.Getenv("AWS_SECRET_KEY"),
		S3Bucket:           os.Getenv("AWS_BUCKET"),
		S3Region:           os.Getenv("AWS_REGION"),
		S3Endpoint:         os.Getenv("AWS_ENDPOINT"),
	}
}
