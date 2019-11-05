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

	// Always print Topic name
	s := []string{
		fmt.Sprintf("Topic : %s", t.Name),
	}

	// Add description if it exists
	if t.Description != "" {
		s = append(s, fmt.Sprintf("Description : %s", t.Description))
	}

	// Add sources if there are any
	if len(t.Sources) > 0 {
		s = append(s, fmt.Sprintf("Sources :"))
		s = append(s, t.Sources...)
	}

	// Always print Last Modified
	s = append(s, fmt.Sprintf("Updated: %s", time.Now().Sub(t.Modified).Truncate(time.Second).String()))

	return strings.Join(s, "\n")
}
