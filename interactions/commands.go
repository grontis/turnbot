package interactions

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type CommandInteraction struct {
	Name        string
	Description string
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (c *CommandInteraction) toCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        c.Name,
		Description: c.Description,
	}
}

type commandManager struct {
	CommandInteractions map[string]*CommandInteraction
}

func newCommandManager() *commandManager {
	return &commandManager{
		CommandInteractions: make(map[string]*CommandInteraction),
	}
}

func (cm *commandManager) registerCommandInteraction(command *CommandInteraction) {
	//Commands don't use CustomID, but are instead accessed by name
	cm.CommandInteractions[command.Name] = command
}

// An exception will be thrown if this is called before the discordgo.Session is opened
func (cm *commandManager) createAllCommands(s *discordgo.Session) error {
	//TODO more dynamically use guildID
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	guildID := os.Getenv("CHANNEL_ID")

	for _, cmd := range cm.CommandInteractions {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildID, cmd.toCommand())
		if err != nil {
			return err
		}
	}
	return nil
}

func (cm *commandManager) handleCommandInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := cm.CommandInteractions[i.ApplicationCommandData().Name]; ok {
		handler.Handler(s, i)
	}
}
