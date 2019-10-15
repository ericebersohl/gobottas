package command

import (
	"fmt"
	"strconv"
)

// Every message that gobottas sees is sent through the parser and the result is a Message
type Message struct{
	// Data provided by the parser
	CommandType CommandType // Enumerated type of the command (if there is no command, type = None)
	Source *Source // Metadata for the message
	Args []string // Any whitespace-delimited arguments given after the command

	// Specific Command Data (provided by interceptors)
	HelpData *HelpData
	MemeData *MemeData
}

// Enum for all defined command types
type CommandType int

const (
	None CommandType = iota
	Error
	Unrecognized
	Help
	Meme
)

func (ct CommandType) String() string {
	return [...]string{"None", "Error", "Unrecognized", "Help", "Meme"}[ct]
}

func StrToCommandType(s string) CommandType {
	switch s {
	case "none":
		return None
	case "help":
		return Help
	case "meme":
		return Meme
	case "error":
		return Error
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

// Defines relevant message data that is persisted after parsing
type Source struct {
	AuthorId Snowflake // Unique identifier of the sender
	ChannelId Snowflake // Unique identifier of the channel in which the message was sent
	Content string // Content of the sent message
}

// Defines the data that the Help command uses
type HelpData struct {
	HelpMsg string // the first line of the help reply; this should always be defined
	SubMsg string // subsequent lines of the help reply; not necessarily defined
}

// Defines the data that the Meme command uses
type MemeData struct {
	Meme string // the meme to be sent
}