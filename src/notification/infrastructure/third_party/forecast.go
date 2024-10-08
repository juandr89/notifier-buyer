package infrastructure

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/juandr89/delivery-notifier-buyer/server"
)

type IForecastService interface {
	FetchForecastByLocation(longitude, latitude, days string) (*ForecastServiceResponse, error)
}

type ForecastService struct {
	BaseURL string
	APIKey  string
}

func NewForecastService(cfg *server.Config) (*ForecastService, error) {
	baseURL := cfg.ForecastServiceConfig.BaseURL
	apiKey := cfg.ForecastServiceConfig.APIKey
	return &ForecastService{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}, nil
}

func (forecast *ForecastService) FetchForecastByLocation(longitude, latitude, days string) (*ForecastServiceResponse, error) {
	url := forecast.BaseURL + "/forecast.json?key=" + forecast.APIKey + "&q=" + latitude + "," + longitude + "&days=" + days + "&aqi=no&alerts=no&lang=es"
	log.Printf("URL: %s", url)

	options := server.RequestOptions{
		Method:         "GET",
		URL:            url,
		MaxRetries:     3,
		RetryDelay:     2 * time.Second,
		RequestTimeout: 5 * time.Second,
	}

	resp, err := server.DoRequestWithRetry(options)
	if err != nil {
		fmt.Println("Request failed:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("FetchForecastByLocation: got http status %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to communicate with the third-party service")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("FetchForecastByLocation: %s", err)
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("FetchDataUbica: %s", err)
		return nil, err
	}

	forecastServiceResponse, err := MapperToDTO(result)

	if err != nil {
		return nil, err
	}

	return forecastServiceResponse, nil
}
