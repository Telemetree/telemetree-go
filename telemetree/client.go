package telemetree

import (
	"github.com/Telemetree/telemetree-go/telemetree/enity"
	"github.com/Telemetree/telemetree-go/telemetree/errors"
	"github.com/Telemetree/telemetree-go/telemetree/internal/config"
	"github.com/Telemetree/telemetree-go/telemetree/internal/encrypt"
	"github.com/Telemetree/telemetree-go/telemetree/internal/rest"
)

type Event = enity.Event

// Client interface defines methods that the Telemetree client should implement.
// The Track method is used to send event data.
type Client interface {
	Track(event Event) error
}

type client struct {
	config *telemetree.Config
	rest   *rest.RestClient
}

// NewClient creates a new instance of the Telemetree client with the provided
// project ID and API key. It initializes the rest client and loads the configuration.
// If there is an error during initialization, a ClientInitializationError is returned.
func NewClient(
	projectID string,
	apiKey string,
) (Client, error) {
	restClient := rest.NewRestClient(apiKey, projectID)
	config, err := telemetree.LoadConfig(apiKey, projectID, restClient)
	if err != nil {
		return nil, errors.NewClientInitializationError(err.Error())
	}

	return &client{
		config: config,
		rest:   restClient,
	}, nil
}

// Track sends an event to the Telemetree API after validating and encrypting it.
// It returns an error if validation, encryption, or sending the event fails.
func (c *client) Track(event Event) error {
	if err := event.Validate(); err != nil {
		return err
	}

	encryptEvent, err := encrypt.PrepareEncryptedPayload(
		c.config.PublicKey,
		c.config.APIKey,
		event,
	)

	if err != nil {
		return errors.NewEventPreparationError(err.Error())
	}

	err = c.rest.SendEvent(c.config.ApiHost, encryptEvent)
	if err != nil {
		return errors.NewEventSendError(err.Error())
	}

	return nil
}
