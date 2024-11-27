![](https://tc-images-api.s3.eu-central-1.amazonaws.com/gif_cropped.gif)
# Telemetree Go SDK

The Telemetree Python SDK provides a convenient way to track and analyze Telegram events using the Telemetree service. With this SDK, you can easily capture and send Telegram events to the Telemetree platform for further analysis and insights.

![Alt](https://repobeats.axiom.co/api/embed/18ee5bb9c80b65e0e060cd5b16802b38262b2a87.svg "Repobeats analytics image")

### Features

- Encrypt event data using a hybrid approach with RSA and AES encryption
- Customize the events and commands to track
- Simple and intuitive API for easy integration

### Installation

Install analytics-go using go get:

```shell
go get github.com/TONSolutions/telemetree-go
```

### Usage

Connect the client to your webhook, or pass the event data directly:

```go
package main

import (
	"fmt"
	"github.com/TONSolutions/telemetree-go"
)

func main() {
	// Create a Telemetree SDK client by providing the Project ID and API Key
	client, err := telemetree.NewClient(
		"YOUR_PROJECT_ID",       // Unique identifier for your project
		"YOUR_API_KEY",          // API key for authentication
	)

	if err != nil {
		fmt.Println("Error creating client:", err)
		return
	}

	// Send an event with user and event data
	err = client.Track(telemetree.Event{
		TelegramID: 112294972,     // User's Telegram ID (required)
		EventType:  "web",  // Event type (required)
		IsPremium:  true,          // Premium status flag (optional)

		// The following fields are optional:
		Username:     "username",  // Username
		Firstname:    "firstname", // First name
		Lastname:     "Lastname",  // Last name
		Language:     "en",        // User's language
		ReferrerType: "web",       // Traffic source type
		Referrer:     "0",         // Traffic source
	})

	if err != nil {
		fmt.Println("Error sending event:", err)
		return
	}

	fmt.Println("Event sent successfully")
}
```

### Encryption

The SDK uses RSA encryption to secure event data before sending it to the Telemetree service, ensuring data privacy. The `publicKey` is fetched automatically from the Telemetree configuration service during initialization, so there’s no need to manually set it.

## Other SDKs
Telemetree SDKs are available for various frameworks and environments, making it easy to incorporate powerful analytics into any Telegram Mini App.
- React SDK: https://github.com/TONSolutions/telemetree-react
- Javascript integration: https://github.com/TONSolutions/telemetree-pixel
- Python SDK: https://github.com/TONSolutions/telemetree-python
- .NET SDK: https://github.com/MANABbl4/Telemetree.Net (community-supported)

### License

This SDK is licensed under the MIT License.
### Support

If you have any questions or need assistance, please contact our support team at support@ton.solutions.
