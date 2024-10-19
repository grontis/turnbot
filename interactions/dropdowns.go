package interactions

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DropdownInteraction struct {
	CustomID    string
	Placeholder string
	Options     []discordgo.SelectMenuOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (di *DropdownInteraction) ToDropdown() discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    di.CustomID,
		Placeholder: di.Placeholder,
		Options:     di.Options,
	}
}

type DropdownManager struct {
	DropdownHandlers map[string]*DropdownInteraction
}

func NewDropdownManager() *DropdownManager {
	return &DropdownManager{
		DropdownHandlers: make(map[string]*DropdownInteraction),
	}
}

func (dm *DropdownManager) RegisterDropdownInteraction(dropdown *DropdownInteraction) {
	dm.DropdownHandlers[dropdown.CustomID] = dropdown
}

func (dm *DropdownManager) HandleDropdownInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := dm.DropdownHandlers[i.MessageComponentData().CustomID]; ok {
		handler.Handler(s, i)
	}
}

func (dm *DropdownManager) SendDropdownMessage(s *discordgo.Session, channelID, customID string, content string) error {
	dropdown := dm.GetDropdownByCustomID(customID)
	if dropdown == nil {
		return fmt.Errorf("dropdown with custom ID '%s' not found", customID)
	}

	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					dropdown.ToDropdown(),
				},
			},
		},
	})
	return err
}

func (dm *DropdownManager) GetDropdownByCustomID(customID string) *DropdownInteraction {
	if dropdown, ok := dm.DropdownHandlers[customID]; ok {
		return dropdown
	}
	return nil
}
