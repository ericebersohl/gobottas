package command

import (
	"math/rand"
	"time"
)

type Interceptor func(*Message) *Message

// Adds relevant data to the HelpData struct.  The Executor might look for this data if it runs into an error
func HelpInterceptor(m *Message) *Message {
	data := HelpData{}

	switch m.CommandType {
	case Unrecognized:
		data.HelpMsg = "The command you entered was not understood."
	case Help:
		data.HelpMsg = "Todo: make this helpful."
	case Meme:
		data.HelpMsg = "Use &meme to check if gobottas is working."
	case Error:
		data.HelpMsg = "Gobottas ran into an unhandled error."
	default:
		// return the message with a nil HelpData pointer
		return m
	}

	// add the struct to the message
	m.HelpData = &data

	return m
}

// Adds a meme to the Message to be returned by the executor
func MemeInterceptor(m *Message) *Message {
	// if the type isn't meme, don't add a meme
	if m.CommandType == Meme {
		memeSlice := []string{
			"Valtteri, it's James.",
			"When did I do dangerous driving?",
		}

		// get a seed; we don't care about crypto security
		rand.Seed(time.Now().UnixNano())

		// apply a random meme to the Message
		m.MemeData = &MemeData{
			Meme: memeSlice[rand.Intn(len(memeSlice))],
		}
	}

	return m
}