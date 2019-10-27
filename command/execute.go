package command

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"log"
)

// Function that executes the commands defined by the Message struct
type Executor func(*discordgo.Session, *Message) error

// Executes commands of type Help
func HelpExecutor(s *discordgo.Session, msg *Message) error {

	// command has no help struct
	if msg.HelpData == nil {
		return errors.New("help executor found message with nil help struct")
	}

	// get the full help message
	helpMsg := msg.HelpData.HelpMsg + msg.HelpData.SubMsg

	_, err := s.ChannelMessageSend(msg.Source.ChannelId.String(), helpMsg)
	if err != nil {
		log.Printf("HelpExecutor: %v", err)
		return err
	}

	return nil
}

// Function to execute any Meme commands in Messages
func MemeExecutor(s *discordgo.Session, msg *Message) error {

	// error check
	if msg.MemeData == nil {
		return errors.New("meme executor found message with nil meme struct")
	}

	// send the return message
	_, err := s.ChannelMessageSend(msg.Source.ChannelId.String(), msg.MemeData.Meme)
	if err != nil {
		log.Printf("MemeExecutor: %v", err)
	}

	return nil
}
