package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/ericebersohl/gobottas/discussion"
	"log"
	"strconv"
	"strings"
	"time"
)

// Function that executes the commands defined by the Message struct
type Executor func(Session, *Registry, *Message) error

// Interface to discordgo.Session
type Session interface {
	ChannelMessageSend(string, string) (*discordgo.Message, error)
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
	retMsg := ""

	switch msg.QueueData.Command {
	case discussion.QAdd:
		// check number of arguments
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call add with no additional arguments.\n&dq add [name] [description]\nwhere name is required and description is optional"
			break
		}

		// make topic
		top := discussion.Topic{
			Name:     msg.Args[0],
			Created:  time.Now(),
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
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("topic %s was added", msg.Args[0])

	case discussion.QRemove:
		// check number of arguments
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call remove without specifying the topic name\n&dq remove [name]\nwhere name is required"
			break
		}

		// call remove
		err := r.DQueue.Remove(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("topic %s was removed", msg.Args[0])

	case discussion.QNext:
		// get topic
		t, err := r.DQueue.Next()
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = t.String()

	case discussion.QBump:
		// check for arg
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call bump without specifying the topic name\n&dq bump [name]\nwhere the name is required"
			break
		}

		// call bump
		err := r.DQueue.Bump(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("topic %s was bumped", msg.Args[0])

	case discussion.QSkip:
		// check for arg
		if len(msg.Args) < 1 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call skip without specifying a topic name\n&dq skip [name]\nwhere name is required"
			break
		}

		// call skip
		err := r.DQueue.Skip(msg.Args[0])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("topic %s has been skipped", msg.Args[0])

	case discussion.QAttach:
		// check arg count
		if len(msg.Args) < 2 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call attach without 2 arguments\n&dq attach [name] [url]\nwhere name is the topic name and url is the source to attach"
			break
		}

		// call attach
		err := r.DQueue.Attach(msg.Args[0], msg.Args[1])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("attached url %s to topic %s", msg.Args[1], msg.Args[0])

	case discussion.QDetach:
		// check arg count
		if len(msg.Args) < 2 {
			msg.QueueData.Command = discussion.QError
			retMsg = "Cannot call detach without 2 arguments\n&dq attach [name] [src-num]\nwhere name is the topic name and src-num is the number of the source to detach"
			break
		}

		// get num
		num, err := strconv.Atoi(msg.Args[1])
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		// call detach
		err = r.DQueue.Detach(msg.Args[0], num)
		if err != nil {
			msg.QueueData.Command = discussion.QError
			retMsg = err.Error()
			break
		}

		retMsg = fmt.Sprintf("detached source %d from topic %s", num, msg.Args[0])

	case discussion.QList:
		// call list
		tops := r.DQueue.List()

		// get strings
		var strs []string
		for i, t := range tops {

			f := fmt.Sprintf("`%2d` %s", i, t.Name)
			if t.Description != "" {
				f = fmt.Sprintf("%s: %s", f, t.Description)
			}

			strs = append(strs, f)
		}

		if len(strs) > 0 {
			retMsg = strings.Join(strs, "\n")
		} else {
			retMsg = "No topics in queue"
		}

	case discussion.QError:
		msg.CommandType = Error
		retMsg = "improper use of command &dq"
		break
	}

	_, err := s.ChannelMessageSend(msg.Source.ChannelId.String(), retMsg)
	if err != nil {
		log.Printf("QueueExecutor: %v", err)
		return err
	}

	return nil
}
