package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"turnbot/interactions"
	"turnbot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var botToken string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//TODO put token in AWS credentials parameter store
	botToken = os.Getenv("DISCORD_TOKEN")
}

func main() {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	//TODO utilize/learn state type

	cmdManager := interactions.NewCommandManager(dg)
	btnManager := interactions.NewButtonManager()

	cmdManager.RegisterCommand(&interactions.Command{
		Name:        "hello",
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

	cmdManager.RegisterCommand(&interactions.Command{
		Name:        "createrole",
		Description: "Creates a new role in the server",
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			createRole(s, i.GuildID, "Foo", discordgo.PermissionManageMessages, 0xFF5733)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Role created successfully!", //TODO conditional data based on success/fail
				},
			})
		},
	})

	btnManager.RegisterButtonInteraction(&interactions.ButtonInteraction{
		CustomID: "button_dice_roll",
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

	//TODO GetButton by ID function
	// channelID := "1295611869515612222"
	// _, err = dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
	// 	Content: "Click the button to roll a dice:",
	// 	Components: []discordgo.MessageComponent{
	// 		discordgo.ActionsRow{
	// 			Components: btnManager.GetButtons(),
	// 		},
	// 	},
	// })
	// if err != nil {
	// 	fmt.Println("Error sending message:", err)
	// 	return
	// }

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmdManager.HandleCommand(s, i)
		case discordgo.InteractionMessageComponent:
			btnManager.HandleButtonInteraction(s, i)
		}
	})

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	err = cmdManager.RegisterAllCommands()
	if err != nil {
		fmt.Println("Error registering commands:", err)
		return
	}

	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

// TODO role manager?
// TODO need to propagate error out of this so that we can indicated error in discord
func createRole(s *discordgo.Session, guildID string, roleName string, permissions int64, color int) {
	roleParams := &discordgo.RoleParams{
		Name:        roleName,
		Permissions: &permissions,
		Color:       &color,
	}

	role, err := s.GuildRoleCreate(guildID, roleParams)
	if err != nil {
		fmt.Printf("Error creating role: %v\n", err)
		return
	}

	fmt.Printf("Role '%s' created successfully with ID: %s\n", role.Name, role.ID)
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
