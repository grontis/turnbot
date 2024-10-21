package interactions

import (
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type ButtonInteraction struct {
	CustomID identifiers.ButtonCustomID
	Label    string
	Style    discordgo.ButtonStyle
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (bi *ButtonInteraction) toButton() discordgo.Button {
	return discordgo.Button{
		Label:    bi.Label,
		Style:    bi.Style,
		CustomID: string(bi.CustomID),
	}
}

type buttonManager struct {
	ButtonInteractions map[identifiers.ButtonCustomID]*ButtonInteraction
}

func newButtonManager() *buttonManager {
	return &buttonManager{
		ButtonInteractions: make(map[identifiers.ButtonCustomID]*ButtonInteraction),
	}
}

func (bm *buttonManager) registerButtonInteraction(button *ButtonInteraction) {
	bm.ButtonInteractions[button.CustomID] = button
}

func (bm *buttonManager) handleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := bm.ButtonInteractions[identifiers.ButtonCustomID(i.MessageComponentData().CustomID)]; ok {
		handler.Handler(s, i)
	}
}

func (bm *buttonManager) buttonInteraction(customID identifiers.ButtonCustomID) *ButtonInteraction {
	if button, ok := bm.ButtonInteractions[customID]; ok {
		return button
	}
	return nil
}
