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
		idxById:  make(map[string]int),
	}
}

type tracker struct {
	length   int
	messages []*Message
	idxById  map[string]int
}

func (t *tracker) Add(message *Message) error {
	if _, ok := t.idxById[message.ID]; ok {
		return nil
	}

	if len(t.messages) == t.length {
		t.messages = t.messages[1:]
	}

	t.messages = append(t.messages, message)
	t.idxById[message.ID] = len(t.messages) - 1
	return nil
}

func (t *tracker) Delete(id string) error {
	idx, ok := t.idxById[id]
	if !ok {
		return ErrMessageNotFound
	}

	t.messages = append(t.messages[:idx], t.messages[idx+1:]...)
	t.idxById = make(map[string]int)
	for i, msg := range t.messages {
		t.idxById[msg.ID] = i
	}
	return nil
}

func (t *tracker) Message(id string) (*Message, error) {
	idx, ok := t.idxById[id]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return t.messages[idx], nil
}

func (t *tracker) Messages() []*Message {
	return t.messages
}
