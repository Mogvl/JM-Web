package main

import (
	"fmt"
	"os"

	"github.com/Mogvl/JM-Web/server/config"
	"github.com/Mogvl/JM-Web/server/internal/api"
	"github.com/Mogvl/JM-Web/server/internal/database"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)

	cfg := config.Load()

	log.Info("Starting JMComic Web Server...")

	db, err := database.New(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Database init failed: %v", err)
	}

	router := api.NewRouter(cfg, db)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Infof("Server listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
