package interactions

import (
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type ModalInteraction struct {
	CustomID   identifiers.ModalCustomID
	Title      string
	Components []discordgo.MessageComponent
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (mi *ModalInteraction) toModal() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:      mi.Title,
			CustomID:   string(mi.CustomID),
			Components: mi.Components,
		},
	}
}

type modalManager struct {
	ModalHandlers map[identifiers.ModalCustomID]*ModalInteraction
}

func newModalManager() *modalManager {
	return &modalManager{
		ModalHandlers: make(map[identifiers.ModalCustomID]*ModalInteraction),
	}
}

func (mm *modalManager) registerModalInteraction(modal *ModalInteraction) {
	mm.ModalHandlers[modal.CustomID] = modal
}

func (mm *modalManager) handleModalInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := mm.ModalHandlers[identifiers.ModalCustomID(i.ModalSubmitData().CustomID)]; ok {
		handler.Handler(s, i)
	}
}

func (mm *modalManager) modalInteraction(customID identifiers.ModalCustomID) *ModalInteraction {
	return mm.ModalHandlers[customID]
}
