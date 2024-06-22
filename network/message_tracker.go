package network

import (
	"errors"
)

// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker, deleting the oldest message if necessary
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Message returns a message for a given ID. Message is retained in tracker
	Message(id string) (message *Message, err error)
	// Messages returns messages in FIFO order
	Messages() (messages []*Message)
}

// ErrMessageNotFound is an error returned by MessageTracker when a message with specified id is not found
var ErrMessageNotFound = errors.New("message not found")

// ErrInvalidLength is an error returned by NewMessageTracker when length is less than or equal to zero
var ErrInvalidLength = errors.New("length must be greater than zero")

// ErrInvalidMessage is an error returned by MessageTracker when Add() is called with an invalid message
var ErrInvalidMessage = errors.New("invalid message")

// NewMessageTracker creates a new MessageTracker with a fixed length.
func NewMessageTracker(length int) (MessageTracker, error) {
	return newTracker(length)
}
