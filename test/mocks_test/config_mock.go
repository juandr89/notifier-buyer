package mocks

type MockConfig struct {
	Port                  string
	APIKey                string
	NotificationSender    string
	SMTPConfig            MockSMTPConfig
	RedisConfig           MockRedisConfig
	ForecastServiceConfig MockForecastServiceConfig
}

type MockSMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type MockForecastServiceConfig struct {
	BaseURL string
	APIKey  string
}

type MockRedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}
