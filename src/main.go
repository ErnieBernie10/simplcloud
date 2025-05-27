package main

import (
	"log"

	"github.com/ErnieBernie10/simplecloud/src/cmd"
	"github.com/ErnieBernie10/simplecloud/src/internal"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	internal.Load()
	cmd.Execute()
}
