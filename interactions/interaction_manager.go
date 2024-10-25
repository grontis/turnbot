package interactions

import (
	"fmt"
	"turnbot/identifiers"

	"github.com/bwmarrin/discordgo"
)

type InteractionManager struct {
	Session         *discordgo.Session
	buttonManager   *buttonManager
	commandManager  *commandManager
	dropdownManager *dropdownManager
	modalManager    *modalManager
}

func NewInteractionManager(s *discordgo.Session) *InteractionManager {
	return &InteractionManager{
		Session:         s,
		buttonManager:   newButtonManager(),
		commandManager:  newCommandManager(),
		dropdownManager: newDropdownManager(),
		modalManager:    newModalManager(),
	}
}

func (im *InteractionManager) AddInteractionHandler(handler func(s *discordgo.Session, i *discordgo.InteractionCreate)) {
	im.Session.AddHandler(handler)
}

func (im *InteractionManager) AddCommandInteraction(cmd *CommandInteraction) {
	im.commandManager.registerCommandInteraction(cmd)
}

func (im *InteractionManager) HandleCommandInteraction(i *discordgo.InteractionCreate) {
	im.commandManager.handleCommandInteraction(im.Session, i)
}

// Create registered commands in discord. Discord API will throw an error if this is called before the session is opened
func (im *InteractionManager) CreateAllCommands() error {
	err := im.commandManager.createAllCommands(im.Session)
	return err
}

func (im *InteractionManager) AddButtonInteraction(btn *ButtonInteraction) {
	im.buttonManager.registerButtonInteraction(btn)
}

func (im *InteractionManager) HandleButtonInteraction(i *discordgo.InteractionCreate) {
	im.buttonManager.handleButtonInteraction(im.Session, i)
}

func (im *InteractionManager) SendButtonMessage(channelID string, customID identifiers.ButtonCustomID, content string) error {
	if channelID == "" {
		return fmt.Errorf("channelID is empty")
	}

	button := im.buttonManager.buttonInteraction(customID)
	if button == nil {
		return fmt.Errorf("no button found with CustomID: %s", customID)
	}
	err := sendButtonMessage(im.Session, channelID, button, content)
	return err
}

func (im *InteractionManager) AddDropdownInteraction(dropdown *DropdownInteraction) {
	im.dropdownManager.registerDropdownInteraction(dropdown)
}

func (im *InteractionManager) HandleDropdownInteraction(i *discordgo.InteractionCreate) {
	im.dropdownManager.handleDropdownInteraction(im.Session, i)
}

func (im *InteractionManager) SendDropdownMessage(channelID string, customID identifiers.DropdownCustomID, content string) error {
	dropdown := im.dropdownManager.dropdownInteraction(identifiers.DropdownClassSelectCustomID)
	if dropdown == nil {
		return fmt.Errorf("no dropdown found with CustomID: %s", customID)
	}
	err := sendDropdownMessage(im.Session, channelID, dropdown, content)
	return err
}

func (im *InteractionManager) AddModalInteraction(modal *ModalInteraction) {
	im.modalManager.registerModalInteraction(modal)
}

func (im *InteractionManager) HandleModalInteraction(i *discordgo.InteractionCreate) {
	im.modalManager.handleModalInteraction(im.Session, i)
}

func (im *InteractionManager) ModalInteractionResponse(customID identifiers.ModalCustomID) *discordgo.InteractionResponse {
	modalInteraction := im.modalManager.modalInteraction(customID)
	return modalInteraction.toModal()
}

func (im *InteractionManager) SendImage(channelID string, filepath string) error {
	err := sendImage(im.Session, channelID, filepath)
	return err
}
