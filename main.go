package main

import (
	"log"
	"os"

	"turnbot/game"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

//TODO logging library

var botToken string
var guildID string
var channelID string

func initSession() (*discordgo.Session, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	//TODO put token in AWS credentials parameter store
	botToken = os.Getenv("DISCORD_TOKEN")

	s, err := discordgo.New("Bot " + botToken)
	return s, err
}

func main() {
	s, err := initSession()
	if err != nil {
		log.Printf("error initializing session: %s", err)
		return
	}

	//TODO utilize/learn state type

	engine := game.NewGameEngine(s)
	engine.Run()
}

//TODO design like game engine.
//(at one point should have modular pieces totally separate in another package)
//"game" design separate from core modular pieces

//Game
//DM
//User defined?
//AI generated?

//Players/characters
//Players will control 1 (to many?) characters

//character definitions
//simplified TTRPG mechanics?
//stats
//weapons
//spells
//classes

//character creation flow
//Input basic text information about character (name, age)
//Because discord modals only support textinput, the following flow will have to be a workflow of sending btns
//race btn -> class btn ->
//select race, class, stats, weapons, spells etc.
//might need to store in DB?
//Modals/TextInput might be a good way to input this data in a form

//Turn management system
//Combat turns
//on a players turn they will be presented with buttons to interact with the world?
//nearby targets

//movement
//non-visual options:
//players are told how far away targets/objects are?
//instead of free movement like on TT, players could have a smaller set of predetermined movements:
//ex:
//move in melee range of {Target}
//move away from {Target}
//climb pillar
//swim across river

//visual options (PREFERRED):
//integrate a web application with the chat that displays tokens? (Unlikely: discord doesn't support IFRAME)
//each turn/action/movement generate an image that represents the map & tokens? (probably the best and most practical option)
//can tap into p5js to create images?

//Chat messages
//Out of character
//in character
//@messages to send messages to specific groups
//DM messages to send DMs to specific players
//message overlay and structure that shows character speaking (Character picture, Name, Class)
//How to enforce users selecting either in character or ooc messages?

//need to be able to conditionally display buttons/accept interactions
//exs:
//Player should not see attack/movement buttons if in combat and not their turn
//Player should only see buttons/skills/actions related to their class

//Custom emojis for the game

//the party
//Voting system for party decisions
