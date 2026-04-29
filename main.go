package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aiya594/doctor-service/internal/app"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	port := os.Getenv("PORT")
	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run server in goroutine
	go func() {
		if err := app.RunServer(port); err != nil {
			log.Fatal("server error:", err)
		}
	}()

	log.Println("server started on port", port)

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop gRPC server gracefully
	done := make(chan struct{})
	go func() {
		app.Stop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Println("shutdown timeout exceeded, forcing stop")
	}

	app.Close()

	log.Println("application shutdown complete")

}
