package infrastructure

type ForecastServiceResponse struct {
	Code        float64 `json:"code"`
	Description string  `json:"description"`
}
