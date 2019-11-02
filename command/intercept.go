package command

import (
	"github.com/ericebersohl/gobottas/discussion"
	"math/rand"
	"strings"
	"time"
)

type Interceptor func(*Message) error

// Adds relevant data to the HelpData struct.  The Executor might look for this data if it runs into an error
// If the first arg has a defined submessage, add it.
// Todo(ee): load help messages from file on startup?
func HelpInterceptor(m *Message) error {
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
		return nil
	}

	// check if there are any args
	if len(m.Args) > 0 {
		switch strings.ToLower(m.Args[0]) {
		case "meme":
			data.SubMsg = "todo: add better help for the meme"
		default:
			// do nothing
		}
	}

	// add the struct to the message
	m.HelpData = &data

	return nil
}

// Adds a meme to the Message to be returned by the executor
func MemeInterceptor(m *Message) error {
	// if the type isn't meme, don't add a meme
	if m.CommandType == Meme {
		memeSlice := []string{
			"Valtteri, it's James.",
			"When did I do dangerous driving?",
			"Steering wheel! Give me the steering wheel. Hey! Hey! Steering wheel, somebody tell him to give it to me! Come on! Move!",
			"Stay out. IN! IN! IN! IN! IN!",
		}

		// get a seed; we don't care about crypto security
		rand.Seed(time.Now().UnixNano())

		// apply a random meme to the Message
		m.MemeData = &MemeData{
			Meme: memeSlice[rand.Intn(len(memeSlice))],
		}
	}

	return nil
}

// Adds a QueueData struct to the message
func QueueInterceptor(m *Message) error {
	if m.CommandType == Queue {
		data := discussion.QueueData{
			Command: discussion.ArgToQueueCommand(m.Args[0]),
		}

		// remove arg if a valid queue command
		if data.Command != discussion.QError {
			m.Args = m.Args[1:]
		}

		m.QueueData = &data
	}

	return nil
}
