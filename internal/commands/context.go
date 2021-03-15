package commands

import "github.com/bwmarrin/discordgo"

type Context struct {
	Session *discordgo.Session
	Message *discordgo.Message // Both of these so I can do "ctx.Message" and "ctx.Session" without having to import discordgo all the time
	Args    []string           // Arguments of the command
	Handler *CommandHandler    // Handle commands
}
