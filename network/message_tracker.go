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

// NewMessageTracker creates a new MessageTracker with a fixed length.
func NewMessageTracker(length int) MessageTracker {
	return &tracker{
		length:   length,
		messages: make([]*Message, 0),
	}
}

type tracker struct {
	length   int
	messages []*Message
}

func (t *tracker) Add(message *Message) error {
	for _, m := range t.messages {
		if m.ID == message.ID {
			return nil
		}
	}

	if len(t.messages) == t.length {
		t.messages = t.messages[1:]
	}

	t.messages = append(t.messages, message)
	return nil
}

func (t *tracker) Delete(id string) error {
	m := make([]*Message, 0)
	found := false

	for _, message := range t.messages {
		if message.ID != id {
			m = append(m, message)
		} else {
			found = true
		}
	}
	t.messages = m

	if !found {
		return ErrMessageNotFound
	}
	return nil
}

func (t *tracker) Message(id string) (*Message, error) {
	for _, message := range t.messages {
		if message.ID == id {
			return message, nil
		}
	}
	return nil, ErrMessageNotFound
}

func (t *tracker) Messages() []*Message {
	return t.messages
}
