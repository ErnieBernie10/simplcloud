package main

import (
	"github.com/ErnieBernie10/simplecloud/src/internal/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}
	api.Serve()
}
