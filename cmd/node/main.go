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
	addrFlag := flag.String("addr", "", "own host:port as seen by peers/clients (defaults to localhost:<port>)")
	peersFlag := flag.String("peers", "", "comma-separated list of peer host:port addresses")
	dataDir := flag.String("data-dir", "data", "directory for persisted raft state")
	flag.Parse()

	addr := *addrFlag
	if addr == "" {
		addr = "localhost:" + *port
	}

	var peers []string
	if *peersFlag != "" {
		peers = strings.Split(*peersFlag, ",")
	}

	store := kv.NewStore()

	applyFn := func(cmd raft.Command) {
		switch cmd.Op {
		case "put":
			store.Set(cmd.Key, cmd.Value)
		case "delete":
			store.Delete(cmd.Key)
		}
	}

	node := raft.NewNode(*id, addr, peers, applyFn, *dataDir)
	svc := service.NewKVService(store, node)
	handler := api.NewHandler(svc)

	app := fiber.New()
	api.RegisterRoutes(app, handler)
	raft.RegisterRoutes(app, node)

	node.Start()

	log.Printf("raftkv node %s listening on port :%s", *id, *port)
	if err := app.Listen(":" + *port); err != nil {
		log.Fatal(err)
	}
}
