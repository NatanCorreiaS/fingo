package main

import (
	"context"
	"log"
	"natan/fingo/dbsqlite"
	"natan/fingo/service"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := dbsqlite.CheckAndCreate(); err != nil {
		log.Println(err)
	}

	// Create a context that is canceled on OS signals for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start the monthly adjustment scheduler in the background
	service.StartMonthlyAdjustmentScheduler(ctx)

	log.Printf("Server listening on PORT 8000...")
	mux := RouterMux()

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	// Run the server in a goroutine so we can listen for shutdown signals
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// Block until we receive a shutdown signal
	<-ctx.Done()
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully.")
}
