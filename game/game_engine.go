package game

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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
var channelID string

type GameEngine struct {
	Session                *discordgo.Session
	EventManager           *events.EventManager
	InteractionsInitLoader InteractionsInitLoader
	GuildInitLoader        GuildInitLoader
	InteractionManager     *interactions.InteractionManager
	GuildManager           *guild.GuildManager
	PlayerCharacters       []Character
}

func NewGameEngine(s *discordgo.Session, interactionsInitLoader InteractionsInitLoader, guildInitLoader GuildInitLoader) (*GameEngine, error) {
	//TODO more dynamically use guildIDs
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	guildID = os.Getenv("GUILD_ID")
	channelID = os.Getenv("CHANNEL_ID")

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
		PlayerCharacters:       make([]Character, 0),
	}

	engine.init()

	return engine, nil
}

func (ge *GameEngine) init() {
	//TODO this kind of operation might be a good scenario to define an error type for less repetition?
	//https://go.dev/blog/errors-are-values

	ge.InteractionsInitLoader.LoadButtonInteractions(ge)
	ge.InteractionsInitLoader.LoadCommandInteractions(ge)
	ge.InteractionsInitLoader.LoadDropdownInteractions(ge)
	ge.InteractionsInitLoader.LoadModalInteractions(ge)
	ge.InteractionsInitLoader.LoadInteractionsHandler(ge)

	//TODO does the event manager belong here? singleton instead?
	ge.StartEventListeners()
}

func (ge *GameEngine) Run() {

	err := ge.Session.Open()
	if err != nil {
		fmt.Printf("Error opening connection: %s", err)
		return
	}

	//It appears that commands can only be added to discord after the session has been opened?
	ge.InteractionsInitLoader.CreateAllCommands(ge)

	ge.GuildInitLoader.SetupBotChannels(ge, guildID)

	//TODO struct property to manage various channels for the game

	//TODO character creation workflow
	ge.startCharacterCreation()

	awaitTerminateSignal()
	ge.Session.Close()
}

func (ge *GameEngine) StartEventListeners() {
	characterCreatedChan := make(chan interface{})
	classSelectedChan := make(chan interface{})

	// Subscribe to character creation event
	ge.EventManager.Subscribe(events.EventCharacterInfoSubmitted, characterCreatedChan)
	go func() {
		for event := range characterCreatedChan {
			fmt.Println("Character info Event Received:", event)
			// Handle character creation logic here
		}
	}()

	// Subscribe to class selection event
	ge.EventManager.Subscribe(events.EventCharacterClassSubmitted, classSelectedChan)
	go func() {
		for event := range classSelectedChan {
			fmt.Println("Class Selected Event Received:", event)
			// Handle class selection logic here
		}
	}()
}

func awaitTerminateSignal() {
	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (ge *GameEngine) startCharacterCreation() {
	//TODO channel management (different channels for game) for ex this should be in a #character-creation channel?
	err := ge.InteractionManager.SendButtonMessage(channelID, identifiers.ButtonStartCharacterCreationCustomID, "Create a character!")
	if err != nil {
		fmt.Println("error sending message:", err)
	}
}
