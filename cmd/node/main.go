package main

import (
	"flag"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ihsanguldur/raftkv/internal/api"
	"github.com/ihsanguldur/raftkv/internal/kv"
	"github.com/ihsanguldur/raftkv/internal/raft"
	"github.com/ihsanguldur/raftkv/internal/service"
)

func main() {
	port := flag.String("port", "8080", "HTTP port to listen on")
	id := flag.String("id", "", "unique node id")
	peersFlag := flag.String("peers", "", "comma-separated list of peer host:port addresses")
	flag.Parse()

	var peers []string
	if *peersFlag != "" {
		peers = strings.Split(*peersFlag, ",")
	}

	store := kv.NewStore()
	svc := service.NewKVService(store)
	handler := api.NewHandler(svc)
	node := raft.NewNode(*id, peers)

	app := fiber.New()
	api.RegisterRoutes(app, handler)
	raft.RegisterRoutes(app, node)

	node.Start()

	log.Printf("raftkv node listenig on port :%s", *port)
	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
