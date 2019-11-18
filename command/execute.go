package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ericebersohl/gobottas/discord"
	"github.com/ericebersohl/gobottas/discussion"
	"log"
	"strconv"
	"time"
)

// Function that executes the commands defined by the Message struct
type Executor func(Session, *Registry, *Message) error

// Interface to discordgo.Session
type Session interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
	ChannelMessageSendEmbed(string, *discordgo.MessageEmbed) (*discordgo.Message, error)
}

// Executes commands of type Help
func HelpExecutor(s Session, r *Registry, msg *Message) error {

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
func MemeExecutor(s Session, r *Registry, msg *Message) error {

	// error check
	if msg.MemeData == nil {
		return errors.New("meme executor found message with nil meme struct")
	}

	// send the return message
	_, err := s.ChannelMessageSend(msg.Source.ChannelId.String(), msg.MemeData.Meme)
	if err != nil {
		log.Printf("MemeExecutor: %v", err)
		return err
	}

	return nil
}

// Function to execute any Queue commands
func QueueExecutor(s Session, r *Registry, msg *Message) error {

	// error check
	if msg.QueueData == nil {
		return errors.New("queue executor found message with nil queuedata struct")
	}

	// message to be returned by gobottas
	var retEmbed *discordgo.MessageEmbed
	var dError *discord.Error

	switch msg.QueueData.Command {
	case discussion.QAdd:
		// check number of arguments
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq add` with no additional arguments." +
					"```&dq add [name] [description]```" +
					"Name is required, while a description is optional.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// make topic
		top := discussion.Topic{
			Name:     msg.Args[0],
			Created:  time.Now(),
			CreatedBy: msg.Source.Username,
			Modified: time.Now(),
		}

		// check for description
		if len(msg.Args) > 1 {
			top.Description = msg.Args[1]
		}

		// add the topic
		err := r.DQueue.Add(&top)
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if err, ok := err.(discord.Error); ok {
				retEmbed = err.Embed()
				break
			}
			return err
		}

	case discussion.QRemove:
		// check number of arguments
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq remove` without an additional argument" +
				"```&dq remove [name]```" +
				"Where name is a required argument.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// call remove
		err := r.DQueue.Remove(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				return err
			}
			break
		}

	case discussion.QNext:
		// get topic
		t, err := r.DQueue.Next()
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err
			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				return err
			}
			break
		}

		retEmbed = t.Embed()

	case discussion.QBump:
		// check for arg
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq bump` without specifying which topic to bump." +
				"```&dq bump [name]```" +
				"Where the topic name is required.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// call bump
		err := r.DQueue.Bump(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				return err
			}
			break
		}

	case discussion.QSkip:
		// check for arg
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq skip` without specifying a topic to skip." +
				"```&dq skip [name]```" +
				"Where the topic name is required.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// call skip
		err := r.DQueue.Skip(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				log.Printf("returning not embedding err: %v", err.(discord.Error).Name)
				return err
			}
			break
		}

	case discussion.QAttach:
		// check arg count
		if len(msg.Args) < 2 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq attach` without specifying a name and source." +
				"```&dq attach [name] [source]```" +
				"Where name is the topic name and source is a url (including https://) and both are required.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// call attach
		err := r.DQueue.Attach(msg.Args[0], msg.Args[1])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				return err
			}
			break
		}

	case discussion.QDetach:
		// check arg count
		if len(msg.Args) < 2 {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Missing Argument(s)",
				"Cannot call `&dq detach` without specifying a topic name and the number of the source to be detached." +
				"```&dq detach [name] [number]```" +
				"Where name is the topic name, and number is the arabic numeral (e.g.: '5') indicating the index of the source.")
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// get num
		num, err := strconv.Atoi(msg.Args[1])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			err := discord.NewError("Invalid Argument",
				fmt.Sprintf("The argument %v could not be converted into an integer.", msg.Args[1]))
			msg.QueueData.Err = err
			retEmbed = err.Embed()
			break
		}

		// call detach
		err = r.DQueue.Detach(msg.Args[0], num)
		if err != nil {
			msg.QueueData.Command = discussion.QError
			msg.QueueData.Err = err

			if errors.As(err, dError) {
				retEmbed = dError.Embed()
			} else {
				return err
			}
			break
		}

	case discussion.QList:
		// call list
		tops := r.DQueue.List()

		retEmbed = &discordgo.MessageEmbed{}
		for _, t := range tops {
			retEmbed.Fields = append(retEmbed.Fields, &discordgo.MessageEmbedField{
				Name:   t.Name,
				Value:  fmt.Sprintf("%s\n%s", t.Description, time.Now().Sub(t.Created).Truncate(time.Second)),
			})
		}

	case discussion.QError:
		msg.CommandType = Error
		msg.QueueData.Err = discord.NewError("Executor Found Error",
			"A message came to executor already tagged as QError")
		break
	}

	// send an embed if there is one
	if retEmbed != nil {
		_, err := s.ChannelMessageSendEmbed(msg.Source.ChannelId.String(), retEmbed)
		if err != nil {
			log.Printf("QueueExecutor: %v", err)
			return err
		}
	}

	return nil
}
