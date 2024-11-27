package telemetree

import (
	"github.com/TONSolutions/telemetree-go/telemetree/internal/rest"
)

const configUrl = "https://config.ton.solutions/v1/client/config"

type Config struct {
	APIKey    string
	ProjectID string
	ApiHost   string
	PublicKey string
}

type Settings struct {
	Host      []string `json:"host"`
	PublicKey []string `json:"public_key"`
}

func LoadConfig(apiKey, projectID string, rest *rest.RestClient) (*Config, error) {
	cfg := &Config{
		APIKey:    apiKey,
		ProjectID: projectID,
	}

	settings, err := rest.LoadConfig(configUrl)
	if err != nil {
		return nil, err
	}

	cfg.ApiHost = settings.Host
	cfg.PublicKey = settings.PublicKey

	return cfg, nil
}
