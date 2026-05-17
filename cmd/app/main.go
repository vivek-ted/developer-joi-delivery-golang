package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"joi-delivery-golang/cmd/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	newServer := server.NewServer(ctx)

	go func() {
		if err := newServer.Start(":8001"); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server ", err)
		}
	}()

	log.Println("server started on :8001")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server....")
	shutDownContext, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := newServer.ShutDown(shutDownContext); err != nil {
		log.Fatal("Server stopped...", err)
	}
	log.Println("server exited...")
}
