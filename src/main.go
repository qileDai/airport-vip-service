package main

import (
	"airport-vip-service/src/db"
	"airport-vip-service/src/routes"
	"airport-vip-service/src/seed"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/airport_vip.db"
	}

	database, err := db.InitDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if len(os.Args) > 1 && os.Args[1] == "--seed" {
		seed.RunSeed(database)
		fmt.Println("Seed data inserted successfully")
		return
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	routes.SetupRoutes(e, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	e.Logger.Fatal(e.Start(":" + port))
}
