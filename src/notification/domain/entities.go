package domain

import "time"

type DeliveryLocation struct {
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
}

type Notification struct {
	Email             string           `json:"email"`
	DeliveryLocation  DeliveryLocation `json:"location"`
	ForecastCode      float64          `json:"forecast_code"`
	BuyerNotification bool             `json:"buyer_notification"`
	Created_at        time.Time        `json:"created_at"`
}
