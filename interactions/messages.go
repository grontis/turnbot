package interactions

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

func sendImage(s *discordgo.Session, channelID string, filepath string) error {

	//TODO image file validation

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s.ChannelFileSend(channelID, filepath, file)
	return err
}

func sendButtonMessage(s *discordgo.Session, channelID string, button *ButtonInteraction, content string) error {
	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					button.toButton(),
				},
			},
		},
	})
	return err
}

func sendDropdownMessage(s *discordgo.Session, channelID string, dropdown *DropdownInteraction, content string) error {
	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Content: content,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					dropdown.toDropdown(),
				},
			},
		},
	})
	return err
}
