package discussion

import (
	"fmt"
	"strings"
	"time"
)

// Defines data for a Queue command
type QueueData struct {
	Command QueueCommand
	Topic   *Topic
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
}

// format the topic for printing
func (t *Topic) String() string {
	return fmt.Sprintf("Topic: %s\nDescription: %s\nSources:\n%s\nModified: %s\n",
		t.Name, t.Description, strings.Join(t.Sources, "\n"), t.Modified.String())
}
