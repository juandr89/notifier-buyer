package infrastructure

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/juandr89/delivery-notifier-buyer/server"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	third_party "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/third_party"
	"github.com/juandr89/delivery-notifier-buyer/src/notification/usecases"
)

type NotificationHandler struct {
	NotificationRepository domain.NotificationRepository
	NotificationSender     domain.NotificationSender
	Config                 server.Config
}

func NewNotificationHandler(repo domain.NotificationRepository, sender domain.NotificationSender, cfg server.Config) *NotificationHandler {
	return &NotificationHandler{
		NotificationRepository: repo,
		NotificationSender:     sender,
		Config:                 cfg,
	}
}

func isValidEmail(email string) bool {
	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

func (c *NotificationHandler) NotifyBuyer(w http.ResponseWriter, r *http.Request) {

	var requestDataNotification usecases.RequestDataNotification
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&requestDataNotification); err != nil {
		domain.ErrorResponseF(w, "NotifyBuyer", http.StatusBadRequest, "Invalid JSON data")
		return
	}

	if !isValidEmail(requestDataNotification.Email) {
		domain.ErrorResponseF(w, "NotifyBuyer", http.StatusBadRequest, "Invalid email")
		return
	}

	log.Printf("NotifyBuyer request [%s] %s", requestDataNotification.Email, requestDataNotification.Location)

	forecastService, err := third_party.NewForecastService(&c.Config)
	if err != nil {
		domain.ErrorResponseF(w, "NotifyBuyer", http.StatusInternalServerError, err.Error())
		return
	}
	log.Printf("NotifyBuyer request %s", forecastService)

	result, err := usecases.SendNotification(requestDataNotification, forecastService, c.NotificationRepository, c.NotificationSender)

	if err != nil {
		domain.ErrorResponseF(w, "NotifyBuyer", http.StatusInternalServerError, err.Error())
		return
	}

	jsonResponse, _ := json.Marshal(result)

	log.Printf("NotifyBuyer response %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func (c *NotificationHandler) BuyerNotifications(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	log.Printf("BuyerNotifications request [%s]", email)

	result, err := usecases.GetBuyerNotification(email, c.NotificationRepository)

	if err != nil {
		if notFoundErr, ok := err.(*domain.NotFoundError); ok {
			domain.ErrorResponseF(w, "NotifyBuyer", http.StatusNotFound, notFoundErr.Message)
			return
		}

		domain.ErrorResponseF(w, "NotifyBuyer", http.StatusInternalServerError, "Unexpected error has ocurred")
		return
	}

	jsonResponse, _ := json.Marshal(result)

	log.Printf("BuyerNotifications response %s", jsonResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}
