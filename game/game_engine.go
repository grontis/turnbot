package game

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"turnbot/identifiers"
	"turnbot/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var channelID string

type GameEngine struct {
	Session            *discordgo.Session
	InteractionManager *interactions.InteractionManager
	PlayerCharacters   []Character
}

func NewGameEngine(s *discordgo.Session) *GameEngine {
	engine := &GameEngine{
		Session:            s,
		InteractionManager: interactions.NewInteractionManager(s),
		PlayerCharacters:   make([]Character, 0),
	}

	//TODO more dynamically use channel_ids
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	channelID = os.Getenv("CHANNEL_ID")

	engine.init()

	return engine
}

func (ge *GameEngine) init() {
	//TODO this kind of operation might be a good scenario to define an error type for less repetition?
	//https://go.dev/blog/errors-are-values

	loadButtonInteractions(ge)
	loadCommandInteractions(ge)
	loadDropdownInteractions(ge)
	loadModalInteractions(ge)
	loadInteractionsHandler(ge)
}

func (ge *GameEngine) Run() {

	err := ge.Session.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	//It appears that commands can only be added to discord after the session has been opened?
	//TODO better design for this?
	ge.InteractionManager.CreateAllCommands()

	err = ge.InteractionManager.SendButtonMessage(channelID, identifiers.ButtonDiceRollCustomID, "Roll some dice!")
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}

	awaitTerminateSignal()
	ge.Session.Close()
}

func awaitTerminateSignal() {

	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
