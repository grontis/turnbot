package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"turnbot/commands"
	"turnbot/interactions"
	"turnbot/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

//TODO logging library

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

	channelId := "1296707235774205972"

	//TODO utilize/learn state type

	cmdManager := commands.NewCommandManager(dg)
	btnManager := interactions.NewButtonManager()
	modalManager := interactions.NewModalManager()

	// TODO how to better handle the dependency of buttons on another manager?
	registerButtons(btnManager, modalManager)
	registerCommands(cmdManager)
	registerModals(modalManager)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			cmdManager.HandleCommand(s, i)
		case discordgo.InteractionMessageComponent:
			btnManager.HandleButtonInteraction(s, i)
		case discordgo.InteractionModalSubmit:
			modalManager.HandleModalSubmission(s, i)
		}
	})

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	err = cmdManager.CreateAllCommands()
	if err != nil {
		fmt.Println("Error registering commands:", err)
		return
	}

	btnManager.SendButtonMessage(dg, channelId, "open_modal_button", "open modal")

	fmt.Println("Bot is running. Press CTRL+C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func registerCommands(cmdManager *commands.CommandManager) {
	cmdManager.RegisterCommand(&commands.Command{
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
}

// TODO how to better handle the dependency of buttons on another manager?
func registerButtons(btnManager *interactions.ButtonManager, modalManager *interactions.ModalManager) {
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

	btnManager.RegisterButtonInteraction(&interactions.ButtonInteraction{
		CustomID: "open_modal_button", // CustomID matches the handler
		Label:    "Open Modal",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// This handler will trigger when the button is clicked
			modal := modalManager.GetModalByCustomID("user_info_modal")
			s.InteractionRespond(i.Interaction, modal.ToModal()) // Open the modal
		},
	})

	btnManager.RegisterButtonInteraction(&interactions.ButtonInteraction{
		CustomID: "create_character",
		Label:    "Create character",
		Style:    discordgo.PrimaryButton,
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {

		},
	})
}

func registerModals(modalManager *interactions.ModalManager) {
	//TODO maybe someway to make this cleaner?
	modalManager.RegisterModal(&interactions.ModalInteraction{
		CustomID: "user_info_modal",
		Title:    "Enter Your Info",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "username_input",
						Label:       "Enter your username",
						Style:       discordgo.TextInputShort,
						Placeholder: "Username",
						Required:    true,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "age_input",
						Label:       "Enter your age",
						Style:       discordgo.TextInputShort,
						Placeholder: "Age",
						Required:    true,
					},
				},
			},
		},
		Handler: func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			username := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
			age := i.ModalSubmitData().Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("You entered: Username: %s, Age: %s", username, age),
				},
			})
		},
	})
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
