package usecases

import "time"

type Location struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type RequestDataNotification struct {
	Email    string   `json:"email"`
	Location Location `json:"location"`
}

type RequestGetNotification struct {
	Email string `json:"email"`
}

type NotificationServiceResponse struct {
	ForecastCode        float64 `json:"forecast_code"`
	ForecastDescription string  `json:"forecast_description"`
	BuyerNotification   bool    `json:"buyer_notification"`
}

type NotificationHistoryDetail struct {
	NotificationSendAt time.Time `json:"notification_sent_at"`
	Location           Location  `json:"location"`
	ForecastCode       float64   `json:"forecast_code"`
}

type NotificationHistoryServiceResponse struct {
	History []NotificationHistoryDetail `json:"history"`
}
