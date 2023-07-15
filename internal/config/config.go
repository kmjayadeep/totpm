package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Cfg struct {
	DBConnectionString string
}

func Get() Cfg {
	return Cfg{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
	}
}
