package interactions

import (
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type ModalInteraction struct {
	CustomID   identifiers.CustomID
	Title      string
	Components []discordgo.MessageComponent
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (mi *ModalInteraction) ToModal() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:      mi.Title,
			CustomID:   string(mi.CustomID),
			Components: mi.Components,
		},
	}
}

type ModalManager struct {
	ModalHandlers map[identifiers.CustomID]*ModalInteraction
}

func NewModalManager() *ModalManager {
	return &ModalManager{
		ModalHandlers: make(map[identifiers.CustomID]*ModalInteraction),
	}
}

func (mm *ModalManager) RegisterModal(modal *ModalInteraction) {
	mm.ModalHandlers[modal.CustomID] = modal
}

func (mm *ModalManager) HandleModalSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := mm.ModalHandlers[identifiers.CustomID(i.ModalSubmitData().CustomID)]; ok {
		handler.Handler(s, i)
	}
}

func (mm *ModalManager) GetModalByCustomID(customID identifiers.CustomID) *ModalInteraction {
	return mm.ModalHandlers[customID]
}
