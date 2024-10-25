package botinit

import (
	"fmt"
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

// Load the button interactions defined in this method into the game engine.
func (b *BotInteractionsInitLoader) LoadButtonInteractions(engine *game.GameEngine) {
	engine.InteractionManager.AddButtonInteraction(&interactions.ButtonInteraction{
		CustomID: identifiers.ButtonStartCharacterCreationCustomID,
		Label:    "Start Character Creation",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			engine.InteractionManager.SendButtonMessage(i.ChannelID, identifiers.ButtonOpenCharacterInfoModalCustomID, "Enter character details")
			engine.InteractionManager.SendDropdownMessage(i.ChannelID, identifiers.DropdownClassSelectCustomID, "Select your class")
			//TODO need to design a way to await and event when BOTH of these interactions are processed (go channels?)

			//TODO publish event for character creation started
			//inside of the handler of that, wait until it receives events for all of the other input components?

			//TODO once all character creation flow is completed, send message with an overview of character to be created.
			//with an option to redo creation/edit?

			// engine.EventManager.Publish(events.EventCharacterClassSubmitted, "Wizard") //EXAMPLE publish event

			err := s.ChannelMessageDelete(i.ChannelID, i.Message.ID)
			if err != nil {
				fmt.Println("Error deleting message:", err)
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate, // Acknowledge interaction
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

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You selected: %s", selectedValue),
				},
			})

			//TODO add "re-select" prompt?

			engine.EventManager.Publish(events.EventCharacterClassSubmitted, selectedValue)

			//TODO event out that a character class was submitted for a given user.
			//TODO outside of this define event handler that will catch those kind of events
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

			username := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			age := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You entered: Username: %s, Age: %s", username, age),
				},
			})

			//TODO event out that a character info was submitted for a given user.
			//TODO outside of this define event handler that will catch those kind of events
		},
	})
}
