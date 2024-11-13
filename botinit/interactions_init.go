package botinit

import (
	"fmt"
	"time"
	"turnbot/events"
	"turnbot/game"
	"turnbot/identifiers"
	"turnbot/interactions"
	"turnbot/utils"

	"github.com/bwmarrin/discordgo"
)

type BotInteractionsInitLoader struct{}

// Create registered commands in discord. Discord API will throw an error if this is called before the session is opened.
func (b *BotInteractionsInitLoader) CreateAllCommands(engine *game.GameEngine) {
	engine.InteractionManager.CreateAllCommands()
}

func (b *BotInteractionsInitLoader) LoadButtonInteractions(engine *game.GameEngine) {
	engine.InteractionManager.AddButtonInteraction(&interactions.ButtonInteraction{
		CustomID: identifiers.ButtonStartCharacterCreationCustomID,
		Label:    "Start Character Creation",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := getUserFromInteraction(i)
			if user == nil {
				fmt.Printf("No user found for discord interaction")
				return
			}

			engine.EventManager.Publish(events.Event{
				EventType: events.EventCharacterCreationStarted,
				Data: &events.CharacterCreationStartedData{
					UserID:    user.ID,
					Timestamp: time.Now(),
				},
			})

			category, err := engine.GuildManager.FindCategoryByName("turnbot")
			if err != nil {
				fmt.Printf("Error finding category: %s", err)
				return //TODO return error from this method?
			}

			userCharacterCreateChannelName := fmt.Sprintf("%s-create-character", user.Username)
			userChannel, err := engine.GuildManager.TryCreateChannelUnderCategory(userCharacterCreateChannelName, category.ID)
			if err != nil {
				fmt.Printf("Error creating channel: %s", err)
				return //TODO return error from this method?
			}

			engine.InteractionManager.SendButtonMessage(userChannel.ID, identifiers.ButtonOpenCharacterInfoModalCustomID, "Enter character details")
			engine.InteractionManager.SendDropdownMessage(userChannel.ID, identifiers.DropdownClassSelectCustomID, "Select your class")

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
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
}

// Load the command interactions defined in this method into the game engine.
func (b *BotInteractionsInitLoader) LoadCommandInteractions(engine *game.GameEngine) {
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

// Load the dropdown interactions defined in this method into the game engine.
func (b *BotInteractionsInitLoader) LoadDropdownInteractions(engine *game.GameEngine) {
	engine.InteractionManager.AddDropdownInteraction(&interactions.DropdownInteraction{
		CustomID:    identifiers.DropdownClassSelectCustomID,
		Placeholder: "Select your character's class",

		//TODO tie the options provided in this dropdown with the classes defined in game logic/rules
		//that way a common identifier can be shared a reused, less error prone
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
			selectedValue := i.MessageComponentData().Values[0]

			user := getUserFromInteraction(i)
			engine.EventManager.Publish(events.Event{
				EventType: events.EventCharacterClassSubmitted,
				Data: &events.CharacterClassSubmittedData{
					UserID:    user.ID,
					ClassName: selectedValue,
				},
			})

			err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				fmt.Println("Error deleting message:", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
		},
	})
}

// Load the interactions handling defined in this method into the game engine.
func (b *BotInteractionsInitLoader) LoadInteractionsHandler(engine *game.GameEngine) {
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
				fmt.Printf("Unknown component type")
			}

		case discordgo.InteractionModalSubmit:
			engine.InteractionManager.HandleModalInteraction(i)

		default:
			fmt.Printf("Unknown interaction type")
		}
	})
}

// Load the modal interactions defined in this method into the game engine.
func (b *BotInteractionsInitLoader) LoadModalInteractions(engine *game.GameEngine) {
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
			//TODO TextInput sanitation and validation. Extension method or middleware that can be used for all TextInput?
			//Some kind of wrapper struct?
			name := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			age := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

			user := getUserFromInteraction(i)
			engine.EventManager.Publish(events.Event{
				EventType: events.EventCharacterInfoSubmitted,
				Data: &events.CharacterInfoSubmittedData{
					UserID: user.ID,
					Name:   name,
					Age:    age,
				},
			})

			err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				fmt.Println("Error deleting message:", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
		},
	})
}

func getUserFromInteraction(i *discordgo.InteractionCreate) *discordgo.User {
	//i.User.ID only used for DMs
	if i.User != nil {
		return i.User
	}

	//i.Member.User.ID only used for within guilds
	if i.Member != nil {
		return i.Member.User
	}

	return nil
}
