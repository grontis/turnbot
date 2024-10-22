package game

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"turnbot/guild"
	"turnbot/interactions"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var guildID string

type GameEngine struct {
	Session                *discordgo.Session
	InteractionLogicLoader InteractionLogicLoader
	GuildLogicLoader       GuildLogicLoader
	InteractionManager     *interactions.InteractionManager
	GuildManager           *guild.GuildManager
	PlayerCharacters       []Character
}

func NewGameEngine(s *discordgo.Session, interactionLogicLoader InteractionLogicLoader, guildLogicLoader GuildLogicLoader) (*GameEngine, error) {
	guildManager, err := guild.NewGuildManager(s, guildID)
	if err != nil {
		return nil, err
	}

	engine := &GameEngine{
		Session:                s,
		InteractionLogicLoader: interactionLogicLoader,
		GuildLogicLoader:       guildLogicLoader,
		InteractionManager:     interactions.NewInteractionManager(s),
		GuildManager:           guildManager,
		PlayerCharacters:       make([]Character, 0),
	}

	//TODO more dynamically use guildIDs
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	guildID = os.Getenv("GUILD_ID")

	engine.init()

	return engine, nil
}

func (ge *GameEngine) init() {
	//TODO this kind of operation might be a good scenario to define an error type for less repetition?
	//https://go.dev/blog/errors-are-values

	ge.InteractionLogicLoader.LoadButtonInteractions(ge)
	ge.InteractionLogicLoader.LoadCommandInteractions(ge)
	ge.InteractionLogicLoader.LoadDropdownInteractions(ge)
	ge.InteractionLogicLoader.LoadModalInteractions(ge)
	ge.InteractionLogicLoader.LoadInteractionsHandler(ge)
}

func (ge *GameEngine) Run() {

	err := ge.Session.Open()
	if err != nil {
		fmt.Printf("Error opening connection: %s", err)
		return
	}

	//It appears that commands can only be added to discord after the session has been opened?
	ge.InteractionLogicLoader.CreateAllCommands(ge)

	ge.GuildLogicLoader.SetupBotChannels(ge, guildID)

	awaitTerminateSignal()
	ge.Session.Close()
}

func awaitTerminateSignal() {
	fmt.Printf("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
