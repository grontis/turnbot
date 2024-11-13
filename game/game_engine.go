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

	awaitTerminateSignal()
	ge.Session.Close()
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
