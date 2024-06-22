package network_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ChainSafe/gossamer-go-interview/network"
	"github.com/stretchr/testify/assert"
)

func generateMessage(n int) *network.Message {
	return &network.Message{
		ID:     fmt.Sprintf("someID%d", n),
		PeerID: fmt.Sprintf("somePeerID%d", n),
		Data:   []byte{0, 1, 1},
	}
}

func TestMessageTracker_NewMessageTracker(t *testing.T) {
	t.Run("create with negative length", func(t *testing.T) {
		mt, err := network.NewMessageTracker(-5)
		assert.ErrorIs(t, err, network.ErrInvalidLength)
		assert.Nil(t, mt)
	})

	t.Run("create with zero length", func(t *testing.T) {
		mt, err := network.NewMessageTracker(0)
		assert.ErrorIs(t, err, network.ErrInvalidLength)
		assert.Nil(t, mt)
	})
}

func TestMessageTracker_Add(t *testing.T) {
	t.Run("add nil", func(t *testing.T) {
		mt, err := network.NewMessageTracker(5)
		assert.NoError(t, err)

		err = mt.Add(nil)
		assert.ErrorIs(t, err, network.ErrInvalidMessage)
	})

	t.Run("add message with invalid ID", func(t *testing.T) {
		mt, err := network.NewMessageTracker(5)
		assert.NoError(t, err)

		message := &network.Message{
			ID:     "",
			PeerID: "somePeerID",
			Data:   []byte{0, 1, 1},
		}

		err = mt.Add(message)
		assert.ErrorIs(t, err, network.ErrInvalidMessage)
	})

	t.Run("add, get, then all messages", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)
	})

	t.Run("add, get, then all messages, delete some", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)

		for i := 0; i < length-2; i++ {
			err := mt.Delete(generateMessage(i).ID)
			assert.NoError(t, err)
		}

		msgs = mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(3),
			generateMessage(4),
		}, msgs)

	})

	t.Run("not full, with duplicates", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(length - 2))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})

	t.Run("not full, with duplicates from other peers", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			msg := generateMessage(length - 2)
			msg.PeerID = "somePeerID0"
			err := mt.Add(msg)
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})
}

func TestMessageTracker_Cleanup(t *testing.T) {
	t.Run("overflow and cleanup", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})

	t.Run("overflow and cleanup with duplicate", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		for i := length; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})
}

func TestMessageTracker_Delete(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)
		err = mt.Delete("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
	})
}

func TestMessageTracker_Message(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt, err := network.NewMessageTracker(length)
		assert.NoError(t, err)
		msg, err := mt.Message("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
		assert.Nil(t, msg)
	})
}

var benchInputs = []struct {
	queueLength int
	idCount     int
}{
	{queueLength: 1000, idCount: 1000},
	{queueLength: 100_000, idCount: 1000},
	{queueLength: 1_000_000, idCount: 1000},
}

func BenchmarkMessageTracker_Add(b *testing.B) {
	for _, v := range benchInputs {
		b.Run(fmt.Sprintf("queue_length_%d", v.queueLength), func(b *testing.B) {
			mt, err := network.NewMessageTracker(v.queueLength)
			if err != nil {
				b.Fail()
			}

			for i := 0; i < b.N; i++ {
				err := mt.Add(generateMessage(i))
				if err != nil {
					b.Fail()
				}
			}
		})
	}
}

// BenchmarkMessageTracker_Delete puts queueLength sequentially generated messages into the queue and picks idCount
// random message IDs to delete from the queue inside the benchmark loop. The benchmark does not currently cover
// deleting messages that are not in the queue.
func BenchmarkMessageTracker_Delete(b *testing.B) {
	for _, v := range benchInputs {
		b.Run(fmt.Sprintf("queue_length_%d", v.queueLength), func(b *testing.B) {
			mt := makeFilledTracker(b, v.queueLength)
			messages := mt.Messages()
			ids := make([]string, v.idCount)

			for i := 0; i < len(ids); i++ {
				ids[i] = messages[rand.Intn(v.queueLength)].ID
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = mt.Delete(ids[i%len(ids)])
			}
		})
	}
}

func BenchmarkMessageTracker_Message(b *testing.B) {
	for _, v := range benchInputs {
		b.Run(fmt.Sprintf("queue_length_%d", v.queueLength), func(b *testing.B) {
			mt := makeFilledTracker(b, v.queueLength)
			messages := mt.Messages()
			ids := make([]string, v.queueLength*2)

			for i := 0; i < len(ids); i++ {
				if rand.Intn(2) == 0 {
					ids[i] = messages[rand.Intn(v.queueLength)].ID
				} else {
					ids[i] = fmt.Sprintf("notInQueue%d", i)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = mt.Message(ids[i%len(ids)])
			}
		})
	}
}

func BenchmarkMessageTracker_Messages(b *testing.B) {
	for _, v := range benchInputs {
		b.Run(fmt.Sprintf("queue_length_%d", v.queueLength), func(b *testing.B) {
			mt := makeFilledTracker(b, v.queueLength)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = mt.Messages()
			}
		})
	}
}

func makeFilledTracker(b *testing.B, queueLength int) network.MessageTracker {
	mt, err := network.NewMessageTracker(queueLength)
	if err != nil {
		b.Fail()
	}

	for i := 0; i < queueLength; i++ {
		err := mt.Add(generateMessage(i))
		if err != nil {
			b.Fail()
		}
	}
	return mt
}
