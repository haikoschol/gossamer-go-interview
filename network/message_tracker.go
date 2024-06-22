package network

import (
	"container/list"
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
		messages: list.New(),
		elemById: make(map[string]*list.Element),
	}
}

type tracker struct {
	length   int
	messages *list.List
	elemById map[string]*list.Element
}

func (t *tracker) Add(message *Message) error {
	if _, ok := t.elemById[message.ID]; ok {
		return nil
	}

	if t.messages.Len() == t.length {
		front := t.messages.Front()
		t.messages.Remove(front)
		delete(t.elemById, front.Value.(*Message).ID)
	}

	t.elemById[message.ID] = t.messages.PushBack(message)
	return nil
}

func (t *tracker) Delete(id string) error {
	elem, ok := t.elemById[id]
	if !ok {
		return ErrMessageNotFound
	}

	t.messages.Remove(elem)
	delete(t.elemById, id)
	return nil
}

func (t *tracker) Message(id string) (*Message, error) {
	elem, ok := t.elemById[id]
	if !ok {
		return nil, ErrMessageNotFound
	}

	return elem.Value.(*Message), nil
}

func (t *tracker) Messages() []*Message {
	messages := make([]*Message, len(t.elemById))
	i := 0

	for e := t.messages.Front(); e != nil; e = e.Next() {
		messages[i] = e.Value.(*Message)
		i++
	}

	return messages
}
