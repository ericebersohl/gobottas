package gobottas

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
)

// Enumeration of the basic types of commands that Gobottas supports
type Command int

const (
	None Command = iota
	Error
	Unrecognized
	Help
	Meme
	Queue
)

// Get the string value associated with a command type
func (c Command) String() string {
	return [...]string{"None", "Error", "Unrecognized", "Help", "Meme", "Queue"}[c]
}

// Parse select strings into commands; note that there are several Commands that no string will parse into
func StrToCommand(s string) Command {
	switch s {
	case "help":
		return Help
	case "meme":
		return Meme
	case "dq":
		return Queue
	default:
		return Unrecognized
	}
}

// functions for Discord's unique id system
type Snowflake uint64

func (s Snowflake) String() string {
	return strconv.FormatUint(uint64(s), 10)
}

func ToSnowflake(s string) (Snowflake, error) {
	sf, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		fmt.Printf("ToSnowflake: %v\n", err)
		return 0, err
	}

	return Snowflake(sf), nil
}

// Every message that Gobottas sees is parsed into a Message and sent through the rest of the application
type Message struct {
	// Provided by the Parser
	Command Command  // Basic command type of the message
	Source  *Source  // Data from discord about the message origin
	Args    []string // parsed args (if there are any)
	Help    string

	// Initialized by Parser, Modified by Interceptors
	Response *Response
}

// Data parsed from the original discord message
type Source struct {
	AuthorId  Snowflake // Unique id of sender
	Username  string    // Username (not including the number) of the sender
	ChannelId Snowflake // Unique id of channel
	Content   string    // Original content of the message
}

type Response struct {
	ChannelId Snowflake
	Text      string
	Embed     *discordgo.MessageEmbed
}

// Session interfaces with the discordgo Session struct using only the relevant functions for Gobottas
type Session interface {
	ChannelMessageSend(channelId string, msg string) (*discordgo.Message, error)
	ChannelMessageSendEmbed(channelId string, embed *discordgo.MessageEmbed) (*discordgo.Message, error)
}

type Registry interface {
	Parse(*discordgo.Message) (*Message, error)
	Intercept(*Message) error
	Execute(*Message, Session) error
}

// Interceptors modify the message.  Every interceptor is called on every message,
// most of the time this is a no-op
type Interceptor func(*Message) error
