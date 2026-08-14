package main

import (
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/ihsanguldur/raftkv/internal/api"
	"github.com/ihsanguldur/raftkv/internal/kv"
	"github.com/ihsanguldur/raftkv/internal/service"
)

func main() {
	port := flag.String("port", "8080", "HTTP port to listen on")
	flag.Parse()

	store := kv.NewStore()
	svc := service.NewKVService(store)
	handler := api.NewHandler(svc)

	app := fiber.New()
	api.RegisterRoutes(app, handler)

	log.Printf("raftkv node listenig on port :%s", *port)
	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
