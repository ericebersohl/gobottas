package discord

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

// Embed functions; adaptation of https://gist.github.com/Necroforger/8b0b70b1a69fa7828b8ad6387ebb3835

// Struct that wraps functionality around a discordgo message embed
type Embed struct {
	*discordgo.MessageEmbed
	charTotal int
}

// definitions of embed element limits
// https://discordapp.com/developers/docs/resources/channel#embed-limits
const (
	TitleLimit = 256
	DescLimit = 2048
	FieldLimit = 25
	FieldNameLimit = 256
	FieldValueLimit = 1024
	FooterLimit = 2048
	AuthorNameLimit = 256
	TotalCharLimit = 6000
)

// Returns a pointer to an empty Embed
func NewEmbed() *Embed {
	e := Embed{
		&discordgo.MessageEmbed{},
		0,
	}
	return &e
}

// Add a URL for the embed
func (e *Embed) EmbedURL(s string) *Embed {
	e.URL = s
	return e
}

// Set the embed title
func (e *Embed) EmbedTitle(t string) *Embed {
	if len(t) > TitleLimit {
		t = t[:TitleLimit]
	}

	if e.charTotal + len(t) > TotalCharLimit {
		log.Printf("Embed already at char limit. Returning unmodified.")
		return e
	}

	e.Title = t
	e.charTotal += len(t)
	return e
}

// Set the embed description
func (e *Embed) EmbedDescription(d string) *Embed {
	if len(d) > DescLimit {
		d = d[:DescLimit]
	}

	if e.charTotal + len(d) > TotalCharLimit {
		return e
	}

	e.Description = d
	e.charTotal += len(d)
	return e
}

// Add a field to the embed
func (e *Embed) AddField(name, value string, inline bool) *Embed {
	if len(name) > FieldNameLimit {
		name = name[:FieldNameLimit]
	}

	if len(value) > FieldValueLimit {
		value = value[:FieldValueLimit]
	}

	if len(e.Fields) < FieldLimit {
		if e.charTotal + len(name) + len(value) > TotalCharLimit {
			log.Printf("Embed already at char limit. Returning unmodified.")
			return e
		}

		e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
			Name:   name,
			Value:  value,
			Inline: inline,
		})

		e.charTotal+= len(value)
		e.charTotal+= len(name)
	}
	return e
}

// Add a custom timestamp to the embed footer
func (e *Embed) EmbedTimestamp(t time.Time) *Embed {
	e.Timestamp = t.Format(time.RFC3339)
	return e
}

// Set the color of the embed
func (e *Embed) EmbedColor(c int) *Embed {
	e.Color = c
	return e
}

// Set the footer of the embed
func (e *Embed) EmbedFooter(text, icon, proxy string) *Embed {
	if len(text) > FooterLimit {
		text = text[:FooterLimit]
	}

	if e.charTotal + len(text) > TotalCharLimit {
		log.Printf("Embed already at char limit. Returning unmodified.")
		return e
	}

	f := discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      icon,
		ProxyIconURL: proxy,
	}

	e.Footer = &f
	return e
}