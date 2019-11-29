package discussion

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discord"
	"strconv"
	"strings"
	"time"
)

// Defines data for a Queue command
type QueueData struct {
	Command QueueCommand
	Topic   *Topic
	Err     error
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
	CreatedBy   string // original author username of the topic
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

func Interceptor(q *Queue) gb.Interceptor {
	return func(msg *gb.Message) error {

		// skip if not Queue message
		if msg.Command != gb.Queue {
			return nil
		}

		// error if registry doesn't have a queue
		if q == nil {
			return errors.New("cannot intercept with nil queue")
		}

		// Queue commands are sent back on the channel in which they are received
		msg.Response.ChannelId = msg.Source.ChannelId

		// used by errors.As to determine whether to create an Error embed
		var dErr *discord.Error

		// attempt to parse command
		cmd := ArgToQueueCommand(msg.Args[0])
		fmt.Printf("CMD: %s ARGS: %d\n", cmd.String(), len(msg.Args))
		switch cmd {
		case QAdd:
			// check for required arg
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Args", "`add` requires at least one argument:\n`&dq add [name] [description?]`").Embed()
				return nil
			}

			t := Topic{
				Sources:     nil,
				Modified:    time.Now(),
				Created:     time.Now(),
				CreatedBy:   msg.Source.Username,
			}

			t.Name = msg.Args[1]

			// add description if exists
			if len(msg.Args) > 2 {
				t.Description = msg.Args[2]
			}

			// add to queue
			if err := q.Add(&t); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
				} else {
					return err
				}
			}

			return nil

		case QRemove:
			// check for name arg
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Arguments", "Remove requires a name argument:\n`&dq remove [name]`").Embed()
				return nil
			}

			// call remove
			if err := q.Remove(msg.Args[1]); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
				} else {
					return err
				}
			}

			return nil

		case QNext:
			// call next; get topic
			t, err := q.Next()
			if err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
					return nil
				} else {
					return err
				}
			}

			msg.Response.Embed = t.Embed()
			return nil

		case QBump:
			// check args
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Args", "Bump requires a name argument:\n`&dq bump [name]`").Embed()
				return nil
			}

			// call bump
			if err := q.Bump(msg.Args[1]); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
					return nil
				} else {
					return err
				}
			}

			return nil

		case QSkip:
			// check args
			if len(msg.Args) < 2 {
				msg.Response.Embed = discord.NewError("Too Few Args", "Skip requries a name argument:\n`&dq skip [name]`").Embed()
				return nil
			}

			// call skip
			if err := q.Skip(msg.Args[1]); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
					return nil
				} else {
					return err
				}
			}

			return nil

		case QAttach:
			// check args
			if len(msg.Args) < 3 {
				msg.Response.Embed = discord.NewError("Too Few Args", "Attach requires two additional arguments:\n`&dq attach [name] [url]`").Embed()
			}

			// call attach
			if err := q.Attach(msg.Args[1], msg.Args[2]); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
					return nil
				} else {
					return err
				}
			}

			return nil

		case QDetach:
			// check args
			if len(msg.Args) < 3 {
				msg.Response.Embed = discord.NewError("Too Few Args", "Detach requires two arguments:\n`&dq detach [name] [number]`" +
					"\nWhere name is the topic name, and number is the index of the source url to remove.").Embed()
				return nil
			}

			// convert the number arg to int
			num, err := strconv.Atoi(msg.Args[2])
			if err != nil {
				msg.Response.Embed = discord.NewError("String to Integer Conversion Error", err.Error()).Embed()
				return nil
			}

			// call detach
			if err := q.Detach(msg.Args[1], num); err != nil {
				if errors.As(err, &dErr) {
					msg.Response.Embed = dErr.Embed()
					return nil
				} else {
					return err
				}
			}

			return nil

		case QList:
			// call list
			l := q.List()

			// create the return embed
			e := discord.NewEmbed().
				EmbedColor(4289797).
				EmbedTitle("Topics").
				EmbedTimestamp(q.Modified)

			// add items
			for _, top := range l {
				e = e.AddField(top.Name, top.Description, false)
			}

			msg.Response.Embed = e.MessageEmbed
			return nil

		case QError:
			e := discord.NewEmbed().
				EmbedColor(13632027).
				EmbedTitle("Unrecognized Command").
				EmbedDescription("Gobottas did not recognize your command.")

			msg.Response.Embed = e.MessageEmbed
			return nil
		}

		return errors.New("reached end of function without returning from switch")
	}
}