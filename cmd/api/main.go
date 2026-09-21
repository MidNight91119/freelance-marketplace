package main

import (
	"context"
	"log"
	"os"

	"github.com/MidNight91119/freelance-marketplace/internal/api"
	db "github.com/MidNight91119/freelance-marketplace/internal/db/sqlc"
	"github.com/MidNight91119/freelance-marketplace/internal/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(connPool)

	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	addr := config.ServerAddress
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	err = server.Start(addr)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
