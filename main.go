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
	botToken = os.Getenv("DISCORD_TOKEN")
}

func main() {
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

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
