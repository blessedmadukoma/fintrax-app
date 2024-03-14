package api

import (
	"database/sql"
	db "fintrax/db/sqlc"
	"fintrax/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
	config  *utils.Config
}

func NewServer(envPath string) *Server {

	config := utils.LoadEnvConfig(envPath)

	// Test DB
	// conn, err := sql.Open(config.DBDriver, fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", config.DBDriver, config.DB_USER, config.DB_PASSWORD, config.DB_HOST, config.DB_PORT, config.DB_NAME))

	// Live DB
	conn, err := sql.Open(config.DBDriver, config.DB_SOURCE_LIVE)
	if err != nil {
		panic(fmt.Sprintf("could not connect to db: %v", err))
	}

	q := db.New(conn)

	g := gin.Default()

	return &Server{
		queries: q,
		router:  g,
		config:  &config,
	}

}

func (s *Server) Start(port int) {
	s.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome to fintrax!!!"})
	})

	User{}.router(s)
	Auth{}.router(s)

	s.router.Run(fmt.Sprintf(":%v", port))
}
