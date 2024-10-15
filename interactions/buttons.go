package interactions

import "github.com/bwmarrin/discordgo"

type ButtonInteraction struct {
	CustomID string
	Label    string
	Style    discordgo.ButtonStyle
	Handler  func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (bi *ButtonInteraction) ToButton() discordgo.Button {
	return discordgo.Button{
		Label:    bi.Label,
		Style:    bi.Style,
		CustomID: bi.CustomID,
	}
}

type ButtonManager struct {
	ButtonHandlers map[string]*ButtonInteraction
}

func NewButtonManager() *ButtonManager {
	return &ButtonManager{
		ButtonHandlers: make(map[string]*ButtonInteraction),
	}
}

func (bm *ButtonManager) RegisterButtonInteraction(button *ButtonInteraction) {
	bm.ButtonHandlers[button.CustomID] = button
}

func (bm *ButtonManager) HandleButtonInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := bm.ButtonHandlers[i.MessageComponentData().CustomID]; ok {
		handler.Handler(s, i)
	}
}

//TODO function to get button by map key

func (bm *ButtonManager) GetButtons() []discordgo.MessageComponent {
	var buttons []discordgo.MessageComponent
	for _, button := range bm.ButtonHandlers {
		buttons = append(buttons, button.ToButton())
	}
	return buttons
}
