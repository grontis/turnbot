package game

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"turnbot/events"
	"turnbot/guild"
	"turnbot/identifiers"
	"turnbot/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

//TODO refactor any game related code out of this and rename to just engine?
//botengine package
//botgame package?
//rename botinit to engineinit? or just place within the game folder that defines the game?
//engine could just handle registering the discord interactions and setup

var guildID string

type GameEngine struct {
	Session                *discordgo.Session
	EventManager           *events.EventManager
	InteractionsInitLoader InteractionsInitLoader
	GuildInitLoader        GuildInitLoader
	InteractionManager     *interactions.InteractionManager
	GuildManager           *guild.GuildManager
	CharacterManager       *CharacterManager
}

func NewGameEngine(s *discordgo.Session, interactionsInitLoader InteractionsInitLoader, guildInitLoader GuildInitLoader) (*GameEngine, error) {
	//TODO more dynamically use guildIDs
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	guildID = os.Getenv("GUILD_ID")

	guildManager, err := guild.NewGuildManager(s, guildID)
	if err != nil {
		return nil, err
	}

	engine := &GameEngine{
		Session:                s,
		EventManager:           events.NewEventManager(),
		InteractionsInitLoader: interactionsInitLoader,
		GuildInitLoader:        guildInitLoader,
		InteractionManager:     interactions.NewInteractionManager(s),
		GuildManager:           guildManager,
		CharacterManager:       NewCharacterManager(),
	}

	engine.init()

	return engine, nil
}

func (ge *GameEngine) init() {
	ge.InteractionsInitLoader.LoadButtonInteractions(ge)
	ge.InteractionsInitLoader.LoadCommandInteractions(ge)
	ge.InteractionsInitLoader.LoadDropdownInteractions(ge)
	ge.InteractionsInitLoader.LoadModalInteractions(ge)
	ge.InteractionsInitLoader.LoadInteractionsHandler(ge)

	ge.StartEventListeners()
}

func (ge *GameEngine) Run() {

	err := ge.Session.Open()
	if err != nil {
		fmt.Printf("Error opening connection: %s", err)
		return
	}

	//It appears that slash commands can only be added to discord after the session has been opened?
	ge.InteractionsInitLoader.CreateAllCommands(ge)

	ge.GuildInitLoader.SetupBotChannels(ge, guildID)

	ge.populateGeneralChannel()

	ge.debugStatAssignment()

	awaitTerminateSignal()
	ge.Session.Close()
}

func (ge *GameEngine) debugStatAssignment() {
	ge.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			handleButtonPress(s, i)
		}
	})

	startCharacterCreation(ge.Session, "TODO")
}

var (
	startingPoints  = 15
	stats           = map[string]int{"Strength": 0, "Dexterity": 0, "Intelligence": 0}
	remainingPoints = startingPoints
)

func startCharacterCreation(s *discordgo.Session, channelID string) {
	// Initial embed message
	embed := &discordgo.MessageEmbed{
		Title:       "Character Creation: Stat Allocation",
		Description: fmt.Sprintf("You have **%d points** to allocate. Use the buttons below to adjust your stats.", remainingPoints),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Strength", Value: fmt.Sprintf("%d", stats["Strength"]), Inline: true},
			{Name: "Dexterity", Value: fmt.Sprintf("%d", stats["Dexterity"]), Inline: true},
			{Name: "Intelligence", Value: fmt.Sprintf("%d", stats["Intelligence"]), Inline: true},
		},
	}

	// Message action rows with buttons
	actionRows := []discordgo.ActionsRow{
		// Row for Strength adjustment
		createStatAdjustmentRow("Strength"),
		// Row for Dexterity adjustment
		createStatAdjustmentRow("Dexterity"),
		// Row for Intelligence adjustment
		createStatAdjustmentRow("Intelligence"),
		// Confirm button
		{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Confirm",
					Style:    discordgo.PrimaryButton,
					CustomID: "confirm",
					Disabled: remainingPoints != 0,
				},
			},
		},
	}

	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed:      embed,
		Components: []discordgo.MessageComponent{actionRows[0], actionRows[1], actionRows[2], actionRows[3]},
	})
}

func createStatAdjustmentRow(stat string) discordgo.ActionsRow {
	return discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    fmt.Sprintf("+ %s", stat),
				Style:    discordgo.SuccessButton,
				CustomID: fmt.Sprintf("add_%s", stat),
				Disabled: remainingPoints <= 0,
			},
			discordgo.Button{
				Label:    fmt.Sprintf("- %s", stat),
				Style:    discordgo.DangerButton,
				CustomID: fmt.Sprintf("subtract_%s", stat),
				Disabled: stats[stat] <= 0,
			},
		},
	}
}

func handleButtonPress(s *discordgo.Session, i *discordgo.InteractionCreate) {
	stat := ""
	if len(i.MessageComponentData().CustomID) > 4 {
		stat = i.MessageComponentData().CustomID[4:] // Extract stat name from custom ID
	}

	switch i.MessageComponentData().CustomID[:3] {
	case "add":
		if remainingPoints > 0 {
			stats[stat]++
			remainingPoints--
		}
	case "sub":
		if stats[stat] > 0 {
			stats[stat]--
			remainingPoints++
		}
	case "con":
		// Confirm button logic
		confirmStats(s, i)
		return
	}

	// Update the embed message with new stats
	updateEmbed(s, i)
}

func updateEmbed(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Update the embed with the latest stats and remaining points
	embed := &discordgo.MessageEmbed{
		Title:       "Character Creation: Stat Allocation",
		Description: fmt.Sprintf("You have **%d points** to allocate. Use the buttons below to adjust your stats.", remainingPoints),
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Strength", Value: fmt.Sprintf("%d", stats["Strength"]), Inline: true},
			{Name: "Dexterity", Value: fmt.Sprintf("%d", stats["Dexterity"]), Inline: true},
			{Name: "Intelligence", Value: fmt.Sprintf("%d", stats["Intelligence"]), Inline: true},
		},
	}

	// Update the action rows to enable/disable buttons based on remaining points
	actionRows := []discordgo.ActionsRow{
		createStatAdjustmentRow("Strength"),
		createStatAdjustmentRow("Dexterity"),
		createStatAdjustmentRow("Intelligence"),
		{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Confirm",
					Style:    discordgo.PrimaryButton,
					CustomID: "confirm",
					Disabled: remainingPoints != 0,
				},
			},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: []discordgo.MessageComponent{actionRows[0], actionRows[1], actionRows[2], actionRows[3]},
		},
	})
}

func confirmStats(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Final confirmation message
	response := fmt.Sprintf("Character stats confirmed!\nStrength: %d\nDexterity: %d\nIntelligence: %d",
		stats["Strength"], stats["Dexterity"], stats["Intelligence"])
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// TODO move actual logic into botinit package like other definitions
func (ge *GameEngine) StartEventListeners() {
	ge.EventManager.Subscribe(events.Subscription{
		EventType: events.EventCharacterCreationStarted,
		Handler: func(data interface{}) {
			if eventData, ok := data.(*events.CharacterCreationStartedData); ok {
				fmt.Printf("Character creation started for user: %s at %v\n", eventData.UserID, eventData.Timestamp)
				ge.CharacterManager.AddNewCharacter(eventData.UserID, &Character{})
			} else {
				fmt.Println("Unexpected data type for EventCharacterCreationStarted")
			}
		},
	})

	ge.EventManager.Subscribe(events.Subscription{
		EventType: events.EventCharacterInfoSubmitted,
		Handler: func(data interface{}) {
			if eventData, ok := data.(*events.CharacterInfoSubmittedData); ok {
				fmt.Printf("Character info Event Received for user:%s Name=%s Age=%s\n", eventData.UserID, eventData.Name, eventData.Age)
				age, err := strconv.Atoi(eventData.Age)

				if err != nil {
					fmt.Println("Conversion error:", err)
				}

				err = ge.CharacterManager.UpdateCharacterInfo(eventData.UserID, eventData.Name, age)
				if err != nil {
					fmt.Printf("Error updating character: %s", err)
				}

				ge.EventManager.Publish(events.Event{
					EventType: events.EventCharacterUpdated,
					Data: &events.CharacterUpdatedData{
						UserID:    eventData.UserID,
						Timestamp: time.Now(),
					},
				})
			} else {
				fmt.Println("Unexpected data type for EventCharacterInfoSubmitted")
			}
		},
	})

	ge.EventManager.Subscribe(events.Subscription{
		EventType: events.EventCharacterClassSubmitted,
		Handler: func(data interface{}) {
			if eventData, ok := data.(*events.CharacterClassSubmittedData); ok {
				fmt.Printf("Character Class Selected Event Received: for user:%s Class=%s\n", eventData.UserID, eventData.ClassName)
				err := ge.CharacterManager.UpdateCharacterClass(eventData.UserID, eventData.ClassName)
				if err != nil {
					fmt.Printf("Error updating character: %s", err)
				}

				ge.EventManager.Publish(events.Event{
					EventType: events.EventCharacterUpdated,
					Data: &events.CharacterUpdatedData{
						UserID:    eventData.UserID,
						Timestamp: time.Now(),
					},
				})
			} else {
				fmt.Println("Unexpected data type for EventCharacterInfoSubmitted")
			}
		},
	})

	ge.EventManager.Subscribe(events.Subscription{
		EventType: events.EventCharacterRaceSubmitted,
		Handler: func(data interface{}) {
			if eventData, ok := data.(*events.CharacterRaceSubmittedData); ok {
				fmt.Printf("Character Race Selected Event Received: for user:%s Class=%s\n", eventData.UserID, eventData.RaceName)
				err := ge.CharacterManager.UpdateCharacterRace(eventData.UserID, eventData.RaceName)
				if err != nil {
					fmt.Printf("Error updating character: %s", err)
				}

				ge.EventManager.Publish(events.Event{
					EventType: events.EventCharacterUpdated,
					Data: &events.CharacterUpdatedData{
						UserID:    eventData.UserID,
						Timestamp: time.Now(),
					},
				})
			} else {
				fmt.Println("Unexpected data type for EventCharacterRaceSubmitted")
			}
		},
	})

	ge.EventManager.Subscribe(events.Subscription{
		EventType: events.EventCharacterUpdated,
		Handler: func(data interface{}) {
			if eventData, ok := data.(*events.CharacterUpdatedData); ok {
				fmt.Printf("Character Updated Event Received: for user:%s at %v\n", eventData.UserID, eventData.Timestamp)

				category, err := ge.GuildManager.FindCategoryByName("turnbot")
				if err != nil {
					fmt.Printf("Error finding category: %s", err)
					return //TODO return error from this method?
				}

				user, err := ge.GuildManager.UserByID(eventData.UserID)
				if err != nil {
					fmt.Printf("Error finding user by ID: %s", err)
					return
				}

				userCharacterSheetChannelName := fmt.Sprintf("%s-character-sheet", user.Username)

				characterChannel, err := ge.GuildManager.TryCreateChannelUnderCategory(userCharacterSheetChannelName, category.ID)
				if err != nil {
					fmt.Printf("Error creating channel: %s", err)
				}

				character := ge.CharacterManager.PlayerCharacters[user.ID] //todo encapsulate in method
				if character == nil {
					fmt.Printf("No character found for userID: %s", user.ID)
					return
				}

				//TODO delete other messages
				//TODO make channel readonly

				ge.InteractionManager.SendTextMessage(characterChannel.ID, character.ToMessageContent())
			} else {
				fmt.Println("Unexpected data type for EventCharacterInfoSubmitted")
			}
		},
	})
}

func awaitTerminateSignal() {
	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (ge *GameEngine) populateGeneralChannel() {
	channel, err := ge.GuildManager.FindChannelInCategoryByName("turnbot", "general")
	if err != nil {
		fmt.Printf("Error finding channel: %s", err)
	}

	err = ge.InteractionManager.SendButtonMessage(channel.ID, identifiers.ButtonStartCharacterCreationCustomID, "Create a character!")
	if err != nil {
		fmt.Println("error sending message:", err)
	}
}
