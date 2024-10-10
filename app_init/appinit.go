package app_init

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juandr89/delivery-notifier-buyer/server"
)

func RunServer(cfg *server.Config) {

	log.Printf("Server starting...")

	router := Routes(cfg)
	srv := &http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Port),
		Handler:     router,
		IdleTimeout: 20 * time.Second,
	}

	go func() {
		log.Printf("Server started! port: " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error ocurred while starting server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := srv.Shutdown(ctx); err != nil {
	// 	log.Fatalf("Error has ocurred while shutting down server: %v", err)
	// }

	log.Printf("Server shutted down!")
}
