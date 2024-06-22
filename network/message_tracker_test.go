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

func TestMessageTracker_Add(t *testing.T) {
	t.Run("add, get, then all messages", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)

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
		mt := network.NewMessageTracker(length)
		err := mt.Delete("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
	})
}

func TestMessageTracker_Message(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		msg, err := mt.Message("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
		assert.Nil(t, msg)
	})
}

const SMALL = 1000
const MEDIUM = 100_000
const LARGE = 1_000_000

func Benchmark_Add_small_queue(b *testing.B) {
	mt := network.NewMessageTracker(SMALL)
	benchmarkAdd(b, mt)
}

func Benchmark_Add_medium_queue(b *testing.B) {
	mt := network.NewMessageTracker(MEDIUM)
	benchmarkAdd(b, mt)
}

func Benchmark_Add_large_queue(b *testing.B) {
	mt := network.NewMessageTracker(LARGE)
	benchmarkAdd(b, mt)
}

func benchmarkAdd(b *testing.B, mt network.MessageTracker) {
	for i := 0; i < b.N; i++ {
		err := mt.Add(generateMessage(i))
		if err != nil {
			b.Fail()
		}
	}
}

func Benchmark_Delete_small_queue(b *testing.B) {
	benchmarkDelete(b, SMALL, 1000)
}

func Benchmark_Delete_medium_queue(b *testing.B) {
	benchmarkDelete(b, MEDIUM, 1000)
}

//func Benchmark_Delete_large_queue(b *testing.B) {
//	benchmarkDelete(b, LARGE, 1000)
//}

// benchmarkDelete puts queueLength sequentially generated messages into the queue and picks idCount random message IDs
// to delete from the queue inside the benchmark loop. The benchmark does not currently cover deleting messages that are
// not in the queue.
func benchmarkDelete(b *testing.B, queueLength int, idCount int) {
	mt := makeFilledTracker(b, queueLength)
	messages := mt.Messages()
	ids := make([]string, idCount)

	for i := 0; i < len(ids); i++ {
		ids[i] = messages[rand.Intn(queueLength)].ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mt.Delete(ids[i%len(ids)])
	}
}

func Benchmark_Message_small_queue(b *testing.B) {
	benchmarkMessage(b, SMALL)
}

func Benchmark_Message_medium_queue(b *testing.B) {
	benchmarkMessage(b, MEDIUM)
}

//func Benchmark_Message_large_queue(b *testing.B) {
//	benchmarkMessage(b, LARGE)
//}

func benchmarkMessage(b *testing.B, queueLength int) {
	mt := makeFilledTracker(b, queueLength)
	messages := mt.Messages()
	ids := make([]string, queueLength*2)

	for i := 0; i < len(ids); i++ {
		if rand.Intn(2) == 0 {
			ids[i] = messages[rand.Intn(queueLength)].ID
		} else {
			ids[i] = fmt.Sprintf("notInQueue%d", i)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mt.Message(ids[i%len(ids)])
	}
}

func Benchmark_Messages_small_queue(b *testing.B) {
	benchmarkMessages(b, SMALL)
}

func Benchmark_Messages_medium_queue(b *testing.B) {
	benchmarkMessages(b, MEDIUM)
}

//func Benchmark_Messages_large_queue(b *testing.B) {
//	benchmarkMessages(b, LARGE)
//}

// avoid inlining of MessageTracker.Messages()
var messages []*network.Message

func benchmarkMessages(b *testing.B, queueLength int) {
	mt := makeFilledTracker(b, queueLength)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		messages = mt.Messages()
	}
}

func makeFilledTracker(b *testing.B, queueLength int) network.MessageTracker {
	mt := network.NewMessageTracker(queueLength)
	for i := 0; i < queueLength; i++ {
		err := mt.Add(generateMessage(i))
		if err != nil {
			b.Fail()
		}
	}
	return mt
}
