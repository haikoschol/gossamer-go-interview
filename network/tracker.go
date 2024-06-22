package network

import "container/list"

type tracker struct {
	length   int
	messages *list.List
	elemById map[string]*list.Element
}

func newTracker(length int) (*tracker, error) {
	if length <= 0 {
		return nil, ErrInvalidLength
	}

	return &tracker{
		length:   length,
		messages: list.New(),
		elemById: make(map[string]*list.Element),
	}, nil
}

func (t *tracker) Add(message *Message) error {
	if message == nil || message.ID == "" {
		return ErrInvalidMessage
	}

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
