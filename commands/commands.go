package commands

//TODO refactor into commands package?

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

type CommandManager struct {
	Session             *discordgo.Session
	ConfiguredChannelID string
	Commands            []*Command
}

func NewCommandManager(session *discordgo.Session) *CommandManager {
	return &CommandManager{
		Session:  session,
		Commands: []*Command{},
	}
}

func (cm *CommandManager) RegisterCommand(command *Command) {
	cm.Commands = append(cm.Commands, command)
}

func (cm *CommandManager) CreateAllCommands() error {
	for _, cmd := range cm.Commands {
		_, err := cm.Session.ApplicationCommandCreate(cm.Session.State.User.ID, "", &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *CommandManager) HandleCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler := cm.getHandler(i.ApplicationCommandData().Name); handler != nil {
		handler(s, i)
	}
}

func (cm *CommandManager) getHandler(name string) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	for _, cmd := range cm.Commands {
		if cmd.Name == name {
			return cmd.Handler
		}
	}
	return nil
}
