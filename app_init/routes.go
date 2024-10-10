package app_init

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juandr89/delivery-notifier-buyer/middleware"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure"
	redisRepository "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/repository"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/sender"
)

func Routes(cfg *server.Config) *mux.Router {
	log.Println("Loading routes..")
	notificationSender := NewNotificationSender(cfg)
	notificationRepository := NewNotificationRepository(cfg)
	notificationHandler := infrastructure.NewNotificationHandler(notificationRepository, notificationSender, *cfg)

	authMiddleware := middleware.ApiKeyMiddleware(cfg.APIKey)

	router := mux.NewRouter()
	api := router.PathPrefix("/api/v1").Subrouter()

	router.Use(authMiddleware)

	api.HandleFunc("/notifications", notificationHandler.NotifyBuyer).Methods(http.MethodPost)
	api.HandleFunc("/notifications/{email}", notificationHandler.BuyerNotifications).Methods(http.MethodGet)

	return router
}

func NewNotificationRepository(cfg *server.Config) domain.NotificationRepository {
	return redisRepository.NewNotificationRepository(cfg.RedisConfig)
}

func NewNotificationSender(cfg *server.Config) domain.NotificationSender {
	return sender.NewSmtpClient(cfg.SMTPConfig)
}
