package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"golang.org/x/exp/rand"
)

var (
	botToken string
)

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

	dg.AddHandler(messageCreate)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			handleSlashCommand(s, i)
		}
	})

	dg.AddHandler(buttonClickHandler)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}

	registerSlashCommands(dg)

	channelID := "1295243704457760819"
	_, err = dg.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Roll 1d6 🎲",
						Style:    discordgo.PrimaryButton,
						CustomID: "button_dice_roll",
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error sending message: ", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL+C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func registerSlashCommands(s *discordgo.Session) error {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "", &discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Says hello!",
	})
	if err != nil {
		fmt.Println("Cannot create slash command: ", err)
		return err
	}
	return nil
}

func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "hello" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hello from your bot!",
			},
		})
	}
}

func buttonClickHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the button clicked was the dice roll button
	if i.MessageComponentData().CustomID == "button_dice_roll" {
		// Simulate rolling a 6-sided dice
		seed := uint64(time.Now().UnixNano()) // Convert int64 to uint64 for rand.Seed
		rand.Seed(seed)                       // Seed random number generator
		diceRoll := rand.Intn(6) + 1          // Generate a number between 1 and 6

		// Respond to the button click with the dice roll result
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("🎲 You rolled a %d!", diceRoll),
			},
		})
	}
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "Hello World!")
	}
}
