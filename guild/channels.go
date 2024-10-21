package guild

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type channelKey struct {
	Name       string
	CategoryID string
}

type ChannelManager struct {
	Session  *discordgo.Session
	GuildID  string
	Channels map[channelKey]*discordgo.Channel
	//TODO categories property?
}

func NewChannelManager(s *discordgo.Session, guildID string) (*ChannelManager, error) {
	channels, err := getGuildChannelMap(s, guildID)
	if err != nil {
		return nil, err
	}

	return &ChannelManager{
		Session:  s,
		GuildID:  guildID,
		Channels: channels,
	}, nil
}

func (cm *ChannelManager) CreateChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel := cm.tryGetChannelByName(channelName, categoryID)
	if channel != nil {
		return nil, fmt.Errorf("channel %s already exists under category", channelName)
	}

	channel, err := cm.createChannelUnderCategory(channelName, categoryID)
	if err != nil {
		fmt.Println("Error creating channel:", err)
		return nil, err
	}

	cm.Channels[channelKey{Name: channelName, CategoryID: categoryID}] = channel
	return channel, nil
}

// TODO more generic createComplexChannel function?
func (cm *ChannelManager) createChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel, err := cm.Session.GuildChannelCreateComplex(cm.GuildID, discordgo.GuildChannelCreateData{
		Name:     channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func (cm *ChannelManager) tryGetChannelByName(channelName string, categoryID string) *discordgo.Channel {
	return cm.Channels[channelKey{Name: channelName, CategoryID: categoryID}]
}

func getGuildChannelMap(s *discordgo.Session, guildID string) (map[channelKey]*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	channelMap := make(map[channelKey]*discordgo.Channel)
	for _, channel := range channels {
		channelMap[channelKey{Name: channel.Name, CategoryID: channel.ParentID}] = channel
	}

	return channelMap, nil
}
