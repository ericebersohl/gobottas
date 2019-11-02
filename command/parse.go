package command

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
	"regexp"
	"strconv"
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
	args, err := Tokenize(cmd.Source.Content)
	if err != nil {
		return cmd, err
	}

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

func Tokenize(s string) (tok []string, err error) {

	// parse the string using the dark art of regular expressions
	tok = regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`).FindAllString(s, -1)

	// remove quotes
	for i := range tok {

		// only call on strings with quotes
		if string(tok[i][0]) == "\"" {
			tok[i], err = strconv.Unquote(tok[i])
			if err != nil {
				log.Printf("Tokenize failed to Unquote string: %v", err)
				return nil, err
			}
		}
	}

	return tok, nil
}
