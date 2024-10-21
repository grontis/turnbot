package game

import (
	"fmt"
	"turnbot/identifiers"
	"turnbot/interactions"
	"turnbot/utils"

	"github.com/bwmarrin/discordgo"
)

func loadButtonInteractions(engine *GameEngine) {
	engine.InteractionManager.AddButtonInteraction(&interactions.ButtonInteraction{
		CustomID: identifiers.ButtonDiceRollCustomID,
		Label:    "Roll 1d6 🎲",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			diceRoll := utils.RollDice(6)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("🎲 You rolled a %d!", diceRoll),
				},
			})
		},
	})

	engine.InteractionManager.AddButtonInteraction(&interactions.ButtonInteraction{
		CustomID: identifiers.ButtonOpenCharacterInfoModalCustomID,
		Label:    "Enter character info",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			modalResponse := engine.InteractionManager.ModalInteractionResponse(identifiers.ModalCharacterInfoCustomID)
			s.InteractionRespond(i.Interaction, modalResponse)
		},
	})
}

func loadCommandInteractions(engine *GameEngine) {
	engine.InteractionManager.AddCommandInteraction(&interactions.CommandInteraction{
		Name:        "hello", //TODO identifiers.CommandName type
		Description: "Says hello!",
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hello from your bot!",
				},
			})
		},
	})
}

func loadDropdownInteractions(engine *GameEngine) {
	engine.InteractionManager.AddDropdownInteraction(&interactions.DropdownInteraction{
		CustomID:    identifiers.DropdownClassSelectCustomID,
		Placeholder: "Select your character's class",
		Options: []discordgo.SelectMenuOption{
			{
				Label:       "Fighter",
				Value:       "fighter",
				Description: "Sword, shield, strength, and honor",
			},
			{
				Label:       "Wizard",
				Value:       "wizard",
				Description: "A never ending thirst for knowledge of the arcana",
			},
			{
				Label:       "Rogue",
				Value:       "rogue",
				Description: "Shadows and daggers",
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// Get the selected value from the select menu
			selectedColor := i.MessageComponentData().Values[0]

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You selected: %s", selectedColor),
				},
			})
		},
	})
}

func loadInteractionsHandler(engine *GameEngine) {
	engine.InteractionManager.AddInteractionHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			engine.InteractionManager.HandleCommandInteraction(i)

		case discordgo.InteractionMessageComponent:
			component := i.MessageComponentData().ComponentType
			switch component {
			case discordgo.ButtonComponent:
				engine.InteractionManager.HandleButtonInteraction(i)
			case discordgo.SelectMenuComponent:
				engine.InteractionManager.HandleDropdownInteraction(i)
			default:
				fmt.Println("Unknown component type")
			}

		case discordgo.InteractionModalSubmit:
			engine.InteractionManager.HandleModalInteraction(i)

		default:
			fmt.Println("Unknown interaction type")
		}
	})
}

func loadModalInteractions(engine *GameEngine) {
	engine.InteractionManager.AddModalInteraction(&interactions.ModalInteraction{
		CustomID: identifiers.ModalCharacterInfoCustomID,
		Title:    "Enter Your Character's Info",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    string(identifiers.TextInputCharacterName),
						Label:       "Enter your character's name",
						Style:       discordgo.TextInputShort,
						Placeholder: "Character name",
						Required:    true,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    string(identifiers.TextInputCharacterAge),
						Label:       "Enter your character's age",
						Style:       discordgo.TextInputShort,
						Placeholder: "Age",
						Required:    true,
					},
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			username := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			age := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You entered: Username: %s, Age: %s", username, age),
				},
			})
		},
	})
}
