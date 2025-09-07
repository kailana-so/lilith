package main

import (
	"lilith/internal/cli"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cli.Execute()
}
