package main

import (
	"fmt"
	"log"

	"github.com/abdullah-mobin/ebook-manager/config"
	"github.com/abdullah-mobin/ebook-manager/route"
	"github.com/gofiber/fiber/v2"
)

func main() {

	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment: %v\n", err)
	}
	fmt.Println("Emvironment Loaded")
	err = config.InitDB()
	if err != nil {
		log.Fatalf("Error connecting database: %v\n", err)
	}
	fmt.Println("Database Connected")

	app := fiber.New()
	route.SetRoute(app)

	app.Listen(":" + config.Port)
}
