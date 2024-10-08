package infrastructure

import (
	"fmt"
)

func MapperToDTO(body map[string]interface{}) (*ForecastServiceResponse, error) {

	forecast, ok := body["forecast"].(map[string]interface{})

	if !ok {
		return nil, fmt.Errorf("forecast is not a nested JSON object")
	}

	forecastday, ok := forecast["forecastday"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("forecastday is not a nested JSON object")
	}

	secondElement, ok := forecastday[1].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("day is not a nested JSON object")
	}

	day, ok := secondElement["day"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("condition is not a nested JSON object")
	}

	var forecastServiceResponse ForecastServiceResponse

	if condition, ok := day["condition"].(map[string]interface{}); ok {
		forecastServiceResponse.Code = condition["code"].(float64)
		forecastServiceResponse.Description = condition["text"].(string)
	}

	if forecastServiceResponse.Code == 0 || forecastServiceResponse.Description == "" {
		return nil, fmt.Errorf("forecast response could not be map succesfully")
	}

	return &forecastServiceResponse, nil
}
