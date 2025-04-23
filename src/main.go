package main

import (
	"log"

	"github.com/ErnieBernie10/simplecloud/cmd"
	"github.com/ErnieBernie10/simplecloud/internal"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	internal.Load()
	cmd.Execute()
}
