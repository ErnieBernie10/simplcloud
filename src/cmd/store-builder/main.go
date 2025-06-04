package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	dir := os.Getenv("STORE_DIR")
	appFs := os.DirFS(dir)
	buildStore(appFs, dir)
}
