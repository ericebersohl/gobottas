package discord

import (
	"github.com/bwmarrin/discordgo"
)

// Custom error class for discord errors
type Error struct {
	Name string
	Desc string
}

// Must implement the Error interface
func (d Error) Error() string {
	return d.Desc
}

// Return an embed message for the error
func (d Error) Embed() *discordgo.MessageEmbed {
	m := NewEmbed().
		EmbedTitle(d.Name).
		EmbedDescription(d.Desc).
		EmbedColor(13632027)

	return m.MessageEmbed
}

// make a new discord error
func NewError(name, desc string) Error {
	err := Error{
		Name: name,
		Desc: desc,
	}
	return err
}
