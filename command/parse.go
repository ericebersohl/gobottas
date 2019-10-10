package command

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// Parse takes a discord message and converts it to the command.Message type
// The function is guaranteed to return a Message with a CommandType
func Parse(msg *discordgo.Message, reg *Registry) (cmd *Message, err error) {

	// Default command has type none
	cmd = &Message{
		CommandType: None,
	}

	// check for nil Author
	if msg.Author == nil {
		cmd.CommandType = Error
		return cmd, errors.New("discord message has nil author")
	}

	// check for empty content
	if msg.Content == "" {
		cmd.CommandType = Error
		return cmd, errors.New("discord message has empty content string")
	}

	src := Source{}
	src.Content = msg.Content

	// convert the authorId string to snowflake
	src.AuthorId, err = ToSnowflake(msg.Author.ID)
	if err != nil {
		log.Printf("Parse: %v", err)
		cmd.CommandType = Error
		return cmd, err
	}

	// convert the channelId string to snowflake
	src.ChannelId, err = ToSnowflake(msg.ChannelID)
	if err != nil {
		log.Printf("Parse: %v", err)
		cmd.CommandType = Error
		return cmd, err
	}

	cmd.Source = &src

	// get the command type
	args := strings.Split(cmd.Source.Content, " ")

	// if the first char of the first argument is not the command prefix, then the command type is none
	if args[0][0] != reg.CommandPrefix {
		return cmd, nil
	}

	// never execute commands from a bot, CommandType is None
	if msg.Author.Bot {
		return cmd, nil
	}

	// get the command type from the first argument, if the command string is not recognized it will default to
	// command type Unrecognized
	cmd.CommandType = StrToCommandType(args[0][1:])
	if cmd.CommandType == Unrecognized {
		// return now so that args don't get updated
		return cmd, nil
	}

	cmd.Args = args[1:] // args[0] is the command arg

	return cmd, nil
}
