package interactions

import "github.com/bwmarrin/discordgo"

type ModalInteraction struct {
	CustomID   string
	Title      string
	Components []discordgo.MessageComponent
	Handler    func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

func (mi *ModalInteraction) ToModal() *discordgo.InteractionResponse {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:      mi.Title,
			CustomID:   mi.CustomID,
			Components: mi.Components,
		},
	}
}

type ModalManager struct {
	ModalHandlers map[string]*ModalInteraction
}

func NewModalManager() *ModalManager {
	return &ModalManager{
		ModalHandlers: make(map[string]*ModalInteraction),
	}
}

func (mm *ModalManager) RegisterModal(modal *ModalInteraction) {
	mm.ModalHandlers[modal.CustomID] = modal
}

func (mm *ModalManager) HandleModalSubmission(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := mm.ModalHandlers[i.ModalSubmitData().CustomID]; ok {
		handler.Handler(s, i)
	}
}

func (mm *ModalManager) GetModalByCustomID(customID string) *ModalInteraction {
	return mm.ModalHandlers[customID]
}

// TODO create modal struct
// define modal creation and submit functions
// it can be used to tie into button actions?
// func sendButtonMessage(s *discordgo.Session, channelID string) {
// 	_, err := s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
// 		Content: "Click the button to open the modal and provide your info:",
// 		Components: []discordgo.MessageComponent{
// 			discordgo.ActionsRow{
// 				Components: []discordgo.MessageComponent{
// 					discordgo.Button{
// 						Label:    "Open Form",
// 						Style:    discordgo.PrimaryButton,
// 						CustomID: "open_modal_button",
// 					},
// 				},
// 			},
// 		},
// 	})
// 	if err != nil {
// 		fmt.Println("Error sending message:", err)
// 	}
// }

// func handleButtonClick(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	if i.MessageComponentData().CustomID == "open_modal_button" {
// 		modal := &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseModal,
// 			Data: &discordgo.InteractionResponseData{
// 				Title:    "User Input Form",
// 				CustomID: "user_input_modal",
// 				Components: []discordgo.MessageComponent{
// 					discordgo.ActionsRow{
// 						Components: []discordgo.MessageComponent{
// 							discordgo.TextInput{
// 								CustomID:    "username_input",
// 								Label:       "Enter your username",
// 								Style:       discordgo.TextInputShort,
// 								Placeholder: "Username",
// 								Required:    true,
// 							},
// 						},
// 					},
// 					discordgo.ActionsRow{
// 						Components: []discordgo.MessageComponent{
// 							discordgo.TextInput{
// 								CustomID:    "age_input",
// 								Label:       "Enter your age",
// 								Style:       discordgo.TextInputShort,
// 								Placeholder: "Age",
// 								Required:    true,
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}

// 		err := s.InteractionRespond(i.Interaction, modal)
// 		if err != nil {
// 			fmt.Println("Error sending modal:", err)
// 		}
// 	}
// }

// func handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
// 	if i.ModalSubmitData().CustomID == "user_input_modal" {
// 		username := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
// 		age := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

// 		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: fmt.Sprintf("You entered: Username: %s, Age: %s", username, age),
// 			},
// 		})
// 		if err != nil {
// 			fmt.Println("Error responding to modal submission:", err)
// 		}
// 	}
// }
