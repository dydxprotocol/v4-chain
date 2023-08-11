package msgsender

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestIndexerMessageSenderInMemoryCollector(t *testing.T) {
	v := NewIndexerMessageSenderInMemoryCollector()
	expectedOnchainMessages := make([]Message, 0, 3)
	expectedOffchainMessages := make([]Message, 0, 3)
	i := 0
	for ; i < 3; i++ {
		expectedOnchainMessages = append(expectedOnchainMessages, Message{
			Key:   []byte("onchainKey" + strconv.Itoa(i)),
			Value: []byte("onchainValue" + strconv.Itoa(i)),
		})
	}
	for ; i < 6; i++ {
		expectedOffchainMessages = append(expectedOffchainMessages, Message{
			Key:   []byte("offchainKey" + strconv.Itoa(i)),
			Value: []byte("offchainValue" + strconv.Itoa(i)),
		})
	}

	v.SendOffchainData(Message{Key: []byte("offChainThatIsCleared")})
	v.SendOnchainData(Message{Key: []byte("onChainThatIsCleared")})
	v.Clear()
	for _, msg := range expectedOnchainMessages {
		v.SendOnchainData(msg)
	}
	for _, msg := range expectedOffchainMessages {
		v.SendOffchainData(msg)
	}
	v.Close()
	v.SendOnchainData(Message{Key: []byte("onChainAfterClose")})
	v.SendOffchainData(Message{Key: []byte("offChainAfterClose")})
	require.Equal(t, expectedOnchainMessages, v.GetOnchainMessages())
	require.Equal(t, expectedOffchainMessages, v.GetOffchainMessages())
}
