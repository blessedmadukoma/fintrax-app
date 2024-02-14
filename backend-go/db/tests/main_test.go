package db_test

import (
	"database/sql"
	db "fintrax/db/sqlc"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

// type DB_URL struct {
// 	DB_Driver string
// }

type Config struct {
	DBDriver    string `mapstructure:"DB_DRIVER"`
	DBSource    string `mapstructure:"DB_SOURCE"`
	DB_USER     string `mapstructure:"DB_USER"`
	DB_PASSWORD string `mapstructure:"DB_PASSWORD"`
	DB_HOST     string `mapstructure:"DB_HOST"`
	DB_PORT     string `mapstructure:"DB_PORT"`
	DB_NAME     string `mapstructure:"DB_NAME"`
}

func LoadEnvConfig(path string) (config Config) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Cannot load env:", err)
	}

	config.DBDriver = os.Getenv("DB_DRIVER")
	config.DB_USER = os.Getenv("DB_USER")
	config.DB_PASSWORD = os.Getenv("DB_PASSWORD")
	config.DB_HOST = os.Getenv("DB_HOST")
	config.DB_PORT = os.Getenv("DB_PORT")
	config.DB_NAME = os.Getenv("DB_NAME")

	return config
}

var testQuery *db.Queries

func TestMain(m *testing.M) {
	config := LoadEnvConfig("../../.env")

	fmt.Println("DB INFO:", config.DBDriver, config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	conn, err := sql.Open(os.Getenv("DB_DRIVER"), fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv("DB_DRIVER"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME")))

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	testQuery = db.New(conn)

	os.Exit(m.Run())
}
