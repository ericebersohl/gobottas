package discussion

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

// Defines data for a Queue command
type QueueData struct {
	Command QueueCommand
	Topic   *Topic
	Err error
}

// Enum for commands for discussion queues
type QueueCommand int

const (
	QError QueueCommand = iota
	QAdd
	QRemove
	QNext
	QBump
	QSkip
	QAttach
	QDetach
	QList
)

func (qc QueueCommand) String() string {
	return [...]string{"Error", "Add", "Remove", "Next", "Bump", "Skip", "Attach", "Detach", "List"}[qc]
}

// parse a string arg into a QueueCommand
func ArgToQueueCommand(arg string) QueueCommand {
	switch arg {
	case "add":
		return QAdd
	case "remove":
		return QRemove
	case "next":
		return QNext
	case "bump":
		return QBump
	case "skip":
		return QSkip
	case "attach":
		return QAttach
	case "detach":
		return QDetach
	case "list":
		return QList
	default:
		return QError
	}
}

// Defines data for a discrete discussion topic
type Topic struct {
	Name        string   // the name of the topic
	Description string   // longer description of the topic
	Sources     []string // an optional list of links to source articles
	Modified    time.Time
	Created     time.Time
	CreatedBy	string // original author username of the topic
}

// format the topic for embedding
func (t *Topic) Embed() *discordgo.MessageEmbed {
	msg := discordgo.MessageEmbed{}

	msg.Color = 4289797
	msg.Title = t.Name
	msg.Footer = &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Proposed by %s", t.CreatedBy)}
	msg.Timestamp = t.Created.Format(time.RFC3339)
	msg.Description = strings.Join(append([]string{t.Description + "\n"}, t.Sources...), "\n")

	return &msg
}