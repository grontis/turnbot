package interactions

import (
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type DropdownInteraction struct {
	CustomID    identifiers.DropdownCustomID
	Placeholder string
	Options     []discordgo.SelectMenuOption
	Handler     func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (di *DropdownInteraction) toDropdown() discordgo.SelectMenu {
	return discordgo.SelectMenu{
		CustomID:    string(di.CustomID),
		Placeholder: di.Placeholder,
		Options:     di.Options,
	}
}

type dropdownManager struct {
	DropdownHandlers map[identifiers.DropdownCustomID]*DropdownInteraction
}

func newDropdownManager() *dropdownManager {
	return &dropdownManager{
		DropdownHandlers: make(map[identifiers.DropdownCustomID]*DropdownInteraction),
	}
}

func (dm *dropdownManager) registerDropdownInteraction(dropdown *DropdownInteraction) {
	dm.DropdownHandlers[(dropdown.CustomID)] = dropdown
}

func (dm *dropdownManager) handleDropdownInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := dm.DropdownHandlers[identifiers.DropdownCustomID(i.MessageComponentData().CustomID)]; ok {
		handler.Handler(s, i)
	}
}

func (dm *dropdownManager) dropdownInteraction(customID identifiers.DropdownCustomID) *DropdownInteraction {
	if dropdown, ok := dm.DropdownHandlers[customID]; ok {
		return dropdown
	}
	return nil
}
