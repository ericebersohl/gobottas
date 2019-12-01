package mock

import (
	gb "github.com/ericebersohl/gobottas"
	"github.com/ericebersohl/gobottas/discord"
)

type MessageOpt func(*gb.Message)

func NewMessage(c gb.Command, opts ...MessageOpt) *gb.Message {
	m := gb.Message{
		Command: c,
		Source: &gb.Source{
			AuthorId:  0,
			Username:  "",
			ChannelId: 0,
			Content:   "",
		},
		Response: &gb.Response{
			ChannelId: 0,
			Text:      "",
			Embed:     nil,
		},
	}

	for _, opt := range opts {
		opt(&m)
	}

	return &m
}

func WithSource(aid, cid gb.Snowflake, uname, content string) MessageOpt {
	return func(message *gb.Message) {
		s := gb.Source{
			AuthorId:  aid,
			Username:  uname,
			ChannelId: cid,
			Content:   content,
		}

		message.Source = &s
	}
}

func WithArgs(arg ...string) MessageOpt {
	return func(msg *gb.Message) {
		var args []string
		for _, a := range arg {
			args = append(args, a)
		}

		msg.Args = args
	}
}

func WithHelp(h string) MessageOpt {
	return func(msg *gb.Message) {
		msg.Help = h
	}
}

func WithResponse(cid gb.Snowflake, text string, embed *discord.Embed) MessageOpt {
	return func(msg *gb.Message) {
		r := gb.Response{
			ChannelId: cid,
			Text:      text,
			Embed:     embed.MessageEmbed,
		}

		msg.Response = &r
	}
}
