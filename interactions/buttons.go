package interactions

import (
	"fmt"
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type ButtonInteraction struct {
	CustomID identifiers.CustomID
	Label    string
	Style    discordgo.ButtonStyle
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (bi *ButtonInteraction) ToButton() discordgo.Button {
	return discordgo.Button{
		Label:    bi.Label,
		Style:    bi.Style,
		CustomID: string(bi.CustomID),
	}
}

type ButtonManager struct {
	ButtonHandlers map[identifiers.CustomID]*ButtonInteraction
}

func NewButtonManager() *ButtonManager {
	return &ButtonManager{
		ButtonHandlers: make(map[identifiers.CustomID]*ButtonInteraction),
	}
}

func (bm *ButtonManager) RegisterButtonInteraction(button *ButtonInteraction) {
	bm.ButtonHandlers[button.CustomID] = button
}

func (bm *ButtonManager) HandleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := bm.ButtonHandlers[identifiers.CustomID(i.MessageComponentData().CustomID)]; ok {
		handler.Handler(s, i)
	}
}

func (bm *ButtonManager) SendButtonMessage(s *discordgo.Session, channelID string, customID identifiers.CustomID, content string) error {
	button := bm.GetButtonByCustomID(customID)
	if button == nil {
		return fmt.Errorf("button with custom ID '%s' not found", customID)
	}

	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button.ToButton(),
				},
			},
		},
	})
	return err
}

func (bm *ButtonManager) GetButtonByCustomID(customID identifiers.CustomID) *ButtonInteraction {
	if button, ok := bm.ButtonHandlers[customID]; ok {
		return button
	}
	return nil
}

func (bm *ButtonManager) GetButtons() []discordgo.MessageComponent {
	var buttons []discordgo.MessageComponent
	for _, button := range bm.ButtonHandlers {
		buttons = append(buttons, button.ToButton())
	}
	return buttons
}
