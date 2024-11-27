package errors

import (
	"fmt"
)

type ClientInitializationError struct {
	Reason string
}

func (e *ClientInitializationError) Error() string {
	return fmt.Sprintf("Client initialization error: %s", e.Reason)
}

func NewClientInitializationError(reason string) *ClientInitializationError {
	return &ClientInitializationError{Reason: reason}
}

type EventValidationError struct {
	Field string
}

func (e *EventValidationError) Error() string {
	return fmt.Sprintf("Event validation error: field '%s' invalid format", e.Field)
}

func NewEventValidationError(field string) *EventValidationError {
	return &EventValidationError{Field: field}
}

type EventPreparationError struct {
	Reason string
}

func (e *EventPreparationError) Error() string {
	return fmt.Sprintf("Event preparation error: %s", e.Reason)
}

func NewEventPreparationError(reason string) *EventPreparationError {
	return &EventPreparationError{Reason: reason}
}

type EventSendError struct {
	Reason string
}

func (e *EventSendError) Error() string {
	return fmt.Sprintf("Error sending event: %s", e.Reason)
}

func NewEventSendError(reason string) *EventSendError {
	return &EventSendError{Reason: reason}
}
