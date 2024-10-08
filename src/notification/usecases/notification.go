package usecases

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/juandr89/delivery-notifier-buyer/src/notification/domain"
	third_party "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/third_party"
)

func RequireBuyerNotification(context context.Context, repository domain.NotificationRepository, code float64) (*bool, error) {
	notificationCodes, err := repository.GetNotificationCodes(context)

	if err != nil {
		return nil, err
	}

	buyerNotification := false

	if slices.Contains(notificationCodes, strconv.FormatFloat(code, 'f', -1, 64)) {
		buyerNotification = true
	}

	return &buyerNotification, nil
}

func CreateNotification(requestDataNotification RequestDataNotification, code float64, requireBuyerNotification bool) (domain.Notification, error) {

	notification := domain.Notification{
		Email: requestDataNotification.Email,
		DeliveryLocation: domain.DeliveryLocation{
			Longitude: requestDataNotification.Location.Longitude,
			Latitude:  requestDataNotification.Location.Latitude,
		},
		ForecastCode:      code,
		BuyerNotification: requireBuyerNotification,
		Created_at:        time.Now(),
	}

	return notification, nil
}

func SendNotification(requestDataNotification *RequestDataNotification, forecastService third_party.IForecastService, repository domain.NotificationRepository, sender domain.NotificationSender) (*NotificationServiceResponse, error) {

	data, err := forecastService.FetchForecastByLocation(requestDataNotification.Location.Longitude, requestDataNotification.Location.Latitude, "2")
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	requireBuyerNotification, err := RequireBuyerNotification(ctx, repository, data.Code)
	if err != nil {
		return nil, err
	}

	notification, err := CreateNotification(*requestDataNotification, data.Code, *requireBuyerNotification)
	if err != nil {
		return nil, err
	}

	notificationServiceResponse := NotificationServiceResponse{
		ForecastCode:        data.Code,
		ForecastDescription: data.Description,
		BuyerNotification:   *requireBuyerNotification,
	}

	text := fmt.Sprintf(`Hola! Tenemos programada la entrega de tu paquete para mañana, en la dirección de  entrega esperamos un día con %s y por esta razón es posible que tengamos retrasos. Haremos todo a nuestro alcance para cumplir con tu entrega.`,
		strings.ToLower(data.Description))

	if *requireBuyerNotification {
		sender.Send(requestDataNotification.Email, text)
		err := repository.SaveNotification(ctx, notification)
		if err != nil {
			return nil, err
		}
	}

	return &notificationServiceResponse, nil
}

func GetBuyerNotification(email string, repository domain.NotificationRepository) (*NotificationHistoryServiceResponse, error) {
	ctx := context.Background()
	notifications, err := repository.GetNotifications(ctx, email)
	if err != nil {
		return nil, err
	}

	notificationHistoryResponse := NotificationHistoryServiceResponse{
		History: MapEntitiesToDTOs(notifications),
	}

	return &notificationHistoryResponse, nil
}

func NotificationEntityToDTO(notification domain.Notification) NotificationHistoryDetail {
	return NotificationHistoryDetail{
		NotificationSendAt: notification.Created_at,
		Location: Location{
			Latitude:  notification.DeliveryLocation.Latitude,
			Longitude: notification.DeliveryLocation.Longitude,
		},
		ForecastCode: notification.ForecastCode,
	}
}

func MapEntitiesToDTOs(entities []domain.Notification) []NotificationHistoryDetail {
	dtos := make([]NotificationHistoryDetail, len(entities))
	for i, entity := range entities {
		dtos[i] = NotificationEntityToDTO(entity)
	}
	return dtos
}
