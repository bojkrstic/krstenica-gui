package main

import (
	"context"
	"fmt"
	"krstenica/internal/config"
	"krstenica/internal/handler"
	"krstenica/internal/repository"
	"krstenica/internal/service"
	"krstenica/migrations"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	conf, err := config.Load()
	if err != nil {
		log.Fatal("Config failed to load", err)
	}

	dbFiles := migrations.GetPostgresMigrations()

	dbVersion, err := config.PostgresMigrate(conf.DB.URL, conf.Migration, dbFiles)
	if err != nil {
		log.Fatalf("Auto migration failed: %+v", err)
	}
	log.Printf("Migrations at version %d", dbVersion)

	db, err := repository.InitORM(conf.DB)
	if err != nil {
		log.Fatalf("Database init failed %+v", err)
	}

	repo := repository.NewRepository(db)

	newService := service.NewService(repo, conf)

	newHandler := handler.NewHttpHandler(newService, conf, repo)
	newHandler.Init()

	fmt.Println(ctx) // We need to use this ctx for graceful shutdown.

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		select {
		case <-sigChan:
			log.Println("closing the service shortly...")
			cancelFunc()
			return
		}
	}

}
