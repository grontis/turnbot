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
	InteractionManager     *interactions.InteractionManager
	PlayerCharacters       []Character
}

func NewGameEngine(s *discordgo.Session, interactionLogicLoader InteractionLogicLoader) *GameEngine {
	engine := &GameEngine{
		Session:                s,
		InteractionLogicLoader: interactionLogicLoader,
		InteractionManager:     interactions.NewInteractionManager(s),
		PlayerCharacters:       make([]Character, 0),
	}

	//TODO more dynamically use guildIDs
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	guildID = os.Getenv("GUILD_ID")

	engine.init()

	return engine
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
	//TODO better design for this?
	ge.InteractionManager.CreateAllCommands()

	guildManager, err := guild.NewGuildManager(ge.Session, guildID)
	if err != nil {
		log.Printf("error creating guild manager: %s", err)
	}
	turnbotCategoryName := "turnbot"
	turnbotCategory, err := guildManager.TryCreateCategory(turnbotCategoryName)
	if err != nil {
		log.Printf("error creating category %s: %s", turnbotCategoryName, err)
	}

	channel, err := guildManager.TryCreateChannelUnderCategory("bot-test", turnbotCategory.ID)
	if err != nil {
		fmt.Printf("error creating channel: %s", err)
		return
	}

	fmt.Printf("channel: %s (ID: %s)\n", channel.Name, channel.ID)

	awaitTerminateSignal()
	ge.Session.Close()
}

func awaitTerminateSignal() {
	fmt.Printf("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
