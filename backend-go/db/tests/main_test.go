package db_test

import (
	"database/sql"
	db "fintrax/db/sqlc"
	"fmt"
	"log"
	"os"
	"testing"

	"fintrax/utils"

	_ "github.com/lib/pq"
)

var testQuery *db.Queries

func TestMain(m *testing.M) {
	config := utils.LoadEnvConfig("../../.env")

	// fmt.Println("DB INFO:", config.DBDriver, config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME)

	conn, err := sql.Open(config.DBDriver, fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", config.DBDriver, config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))

	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	testQuery = db.New(conn)

	os.Exit(m.Run())
}
