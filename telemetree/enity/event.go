package enity

import "github.com/TONSolutions/telemetree-go/telemetree/errors"

// Event represents the client event data
type Event struct {
	TelegramID   int
	EventType    string
	IsPremium    bool
	Username     string
	Firstname    string
	Lastname     string
	Language     string
	ReferrerType string
	Referrer     string
}

func (e *Event) Validate() error {
	if e.TelegramID == 0 {
		return errors.NewEventValidationError("TelegramID")
	}
	if e.EventType == "" {
		return errors.NewEventValidationError("EventType")
	}
	return nil
}
