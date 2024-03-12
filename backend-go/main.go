package main

import (
	"fintrax/api"
)

func main() {
	server := api.NewServer(".env")

	server.Start(8081)
}
