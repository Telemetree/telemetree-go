package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	timeoutHTP       = 10 * time.Second
	headerXApiKey    = "x-api-key"
	headerXProjectID = "x-project-id"
)

type RestClient struct {
	client    *http.Client
	apiKey    string
	projectID string
}

func NewRestClient(apiKey, projectID string) *RestClient {
	return &RestClient{
		apiKey:    apiKey,
		projectID: projectID,
		client: &http.Client{
			Timeout: timeoutHTP,
		},
	}
}

type SettingsResponse struct {
	Host      string `json:"host"`
	PublicKey string `json:"public_key"`
}

func (rc *RestClient) buildURL(host, endpoint string, params map[string]string) string {
	var url string

	if endpoint == "" {
		url = host
	} else {
		url = fmt.Sprintf("%s/%s", host, endpoint)
	}

	if params != nil {
		queryParams := "?"
		for key, value := range params {
			queryParams += fmt.Sprintf("%s=%s&", key, value)
		}
		url += queryParams[:len(queryParams)-1]
	}
	return url
}

func (rc *RestClient) createRequest(
	method,
	url string,
	body []byte,
	headers map[string]string,
) (*http.Request, error) {

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rc.apiKey))
	req.Header.Set("Content-Type", "application/json")

	if headers != nil {
		for header, value := range headers {
			req.Header.Set(header, value)
		}
	}

	return req, nil
}

func (rc *RestClient) sendRequest(req *http.Request) (*http.Response, error) {
	resp, err := rc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	return resp, nil
}

func (rc *RestClient) LoadConfig(host string) (SettingsResponse, error) {
	var response SettingsResponse
	url := rc.buildURL(host, "", map[string]string{
		"project": rc.projectID,
	})

	req, err := rc.createRequest(http.MethodGet, url, nil, nil)
	if err != nil {
		return response, err
	}

	resp, err := rc.sendRequest(req)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var settings SettingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return response, fmt.Errorf("failed to decode response: %w", err)
	}

	return settings, nil
}

func (rc *RestClient) SendEvent(host string, event []byte) error {
	url := rc.buildURL(host, "", nil)

	headers := map[string]string{
		headerXApiKey:    rc.apiKey,
		headerXProjectID: rc.projectID,
	}

	req, err := rc.createRequest(http.MethodPost, url, event, headers)
	if err != nil {
		return err
	}

	resp, err := rc.sendRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
