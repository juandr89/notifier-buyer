package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/juandr89/delivery-notifier-buyer/middleware"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure"
	redisRepository "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/repository"
	sender "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/sender"
)

func RunServer(cfg *server.Config) {

	log.Printf("Server starting...")

	notificationSender := createNotificationSender(cfg)
	notificationRepository := createNotificationRepository(cfg)
	notificationHandler := infrastructure.NewNotificationHandler(notificationRepository, notificationSender, *cfg)

	authMiddleware := middleware.ApiKeyMiddleware(cfg.APIKey)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	router.Use(authMiddleware)

	api.HandleFunc("/notifications", notificationHandler.NotifyBuyer).Methods(http.MethodPost)
	api.HandleFunc("/notifications/{email}", notificationHandler.BuyerNotifications).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
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
	log.Printf("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Error has ocurred while shutting down server: %v", err)
	}

	log.Printf("Server shutted down!")
}

func createNotificationRepository(cfg *server.Config) domain.NotificationRepository {
	return redisRepository.NewNotificationRepository(cfg.RedisConfig)
}

func createNotificationSender(cfg *server.Config) domain.NotificationSender {
	return sender.NewSmtpClient(cfg.SMTPConfig)
}

func main() {
	cfg, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	RunServer(cfg)
}
