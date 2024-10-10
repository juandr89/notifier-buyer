package service_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"bou.ke/monkey"
	"github.com/golang/mock/gomock"
	"github.com/juandr89/delivery-notifier-buyer/server"

	third_party "github.com/juandr89/delivery-notifier-buyer/src/notification/infrastructure/third_party"
	"github.com/stretchr/testify/assert"
)

func TestMapperToDTO(t *testing.T) {
	t.Run("Successful mapping", func(t *testing.T) {
		input := map[string]interface{}{
			"forecast": map[string]interface{}{
				"forecastday": []interface{}{
					map[string]interface{}{
						"day": map[string]interface{}{
							"condition": map[string]interface{}{
								"code": 300.0,
								"text": "Partly cloudy",
							},
						},
					},
					map[string]interface{}{
						"day": map[string]interface{}{
							"condition": map[string]interface{}{
								"code": 200.0,
								"text": "Sunny",
							},
						},
					},
				},
			},
		}
		expectedResponse := &third_party.ForecastServiceResponse{
			Code:        200.0,
			Description: "Sunny",
		}

		response, err := third_party.MapperToDTO(input)

		assert.NoError(t, err)
		assert.Equal(t, expectedResponse, response)
	})
	t.Run("MissingForecastKey", func(t *testing.T) {
		input := map[string]interface{}{
			"foo": "bar",
		}
		expectedError := "forecast is not a nested JSON object"

		response, err := third_party.MapperToDTO(input)

		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError)
	})

	t.Run("MissingForecastdayKey", func(t *testing.T) {
		input := map[string]interface{}{
			"forecast": map[string]interface{}{
				"foo": "bar",
			},
		}
		expectedError := "forecastday is not a nested JSON object"

		response, err := third_party.MapperToDTO(input)

		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError)
	})

	t.Run("InvalidForecastdayDataStructure", func(t *testing.T) {
		input := map[string]interface{}{
			"forecast": map[string]interface{}{
				"forecastday": "not-an-array",
			},
		}
		expectedError := "forecastday is not a nested JSON object"

		response, err := third_party.MapperToDTO(input)

		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError)
	})
	t.Run("InvalidDayAsSecondElement", func(t *testing.T) {
		input := map[string]interface{}{
			"forecast": map[string]interface{}{
				"forecastday": []interface{}{
					map[string]interface{}{
						"day": map[string]interface{}{
							"condition": map[string]interface{}{
								"code": 200.0,
								"text": "Sunny",
							},
						},
					},
					"not-a-map",
				},
			},
		}
		expectedError := "day is not a nested JSON object"

		response, err := third_party.MapperToDTO(input)

		assert.Nil(t, response)
		assert.EqualError(t, err, expectedError)
	})
}

func TestFetchForecastByLocation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	forecastService := &third_party.ForecastService{
		BaseURL: "http://example.com",
		APIKey:  "apikey",
	}

	t.Run("Success", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"forecast": "test"}`)),
		}

		mockResponseService := &third_party.ForecastServiceResponse{}
		monkey.Patch(server.DoRequestWithRetry, func(opts server.RequestOptions) (*http.Response, error) {
			return mockResponse, nil
		})
		defer monkey.Unpatch(server.DoRequestWithRetry)

		monkey.Patch(third_party.MapperToDTO, func(body map[string]interface{}) (*third_party.ForecastServiceResponse, error) {
			return mockResponseService, nil
		})
		defer monkey.Unpatch(server.DoRequestWithRetry)

		result, err := forecastService.FetchForecastByLocation("123.456", "78.910", "3")

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("RequestToServiceFailed", func(t *testing.T) {
		monkey.Patch(server.DoRequestWithRetry, func(opts server.RequestOptions) (*http.Response, error) {
			return nil, errors.New("failed to make request")
		})
		defer monkey.Unpatch(server.DoRequestWithRetry)

		result, err := forecastService.FetchForecastByLocation("123.456", "78.910", "3")

		assert.EqualError(t, err, "failed to make request")
		assert.Nil(t, result)
	})

	t.Run("StatusCodeIsNot200", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 500,
			Body:       io.NopCloser(bytes.NewBufferString("")),
		}
		monkey.Patch(server.DoRequestWithRetry, func(opts server.RequestOptions) (*http.Response, error) {
			return mockResponse, nil
		})
		defer monkey.Unpatch(server.DoRequestWithRetry)

		result, err := forecastService.FetchForecastByLocation("123.456", "78.910", "3")

		assert.EqualError(t, err, "failed to communicate with the third-party service")
		assert.Nil(t, result)
	})

	t.Run("InvalidJSONResponse", func(t *testing.T) {
		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{invalid-json}`)),
		}
		monkey.Patch(server.DoRequestWithRetry, func(opts server.RequestOptions) (*http.Response, error) {
			return mockResponse, nil
		})
		defer monkey.Unpatch(server.DoRequestWithRetry)

		result, err := forecastService.FetchForecastByLocation("123.456", "78.910", "3")

		assert.EqualError(t, err, "invalid character 'i' looking for beginning of object key string")
		assert.Nil(t, result)
	})

}
