package guild

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type channelKey struct {
	Name     string
	ParentID string
}

// TODO encapsulate members
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

func (cm *ChannelManager) TryCreateCategory(categoryName string) (*discordgo.Channel, error) {
	category := cm.CategoryByName(categoryName)
	if category != nil {
		return category, nil
	}

	category, err := cm.createCategory(categoryName)
	if err != nil {
		return nil, fmt.Errorf("error creating category: %v", err)
	}

	//TODO categories map?
	cm.Channels[channelKey{Name: categoryName, ParentID: ""}] = category
	return category, nil
}

// Find category channel by name. Returns nil if none found. (discord categories are represented as discordgo.Channel type)
func (cm *ChannelManager) CategoryByName(categoryName string) *discordgo.Channel {
	for _, channel := range cm.Channels {
		if channel.Name == categoryName && channel.Type == discordgo.ChannelTypeGuildCategory {
			return channel
		}
	}

	return nil
}

// First tries to retrieve an existing channel. If not found, will attempt to create channel.
func (cm *ChannelManager) TryCreateChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel := cm.ChannelByName(channelName, categoryID)
	if channel != nil {
		return channel, nil
	}

	channel, err := cm.createChannelUnderCategory(channelName, categoryID)
	if err != nil {
		fmt.Printf("Error creating channel: %s", err)
		return nil, err
	}

	cm.Channels[channelKey{Name: channelName, ParentID: categoryID}] = channel
	return channel, nil
}

func (cm *ChannelManager) ChannelByName(channelName string, categoryID string) *discordgo.Channel {
	return cm.Channels[channelKey{Name: channelName, ParentID: categoryID}]
}

func (cm *ChannelManager) ChannelsUnderCategory(categoryName string) []*discordgo.Channel {
	var result []*discordgo.Channel
	var categoryID string

	for _, channel := range cm.Channels {
		if channel.Name == categoryName && channel.Type == discordgo.ChannelTypeGuildCategory {
			categoryID = channel.ID
			break
		}
	}

	if categoryID == "" {
		fmt.Printf("Category '%s' not found\n", categoryName)
		return result
	}

	for _, channel := range cm.Channels {
		if channel.ParentID == categoryID {
			result = append(result, channel)
		}
	}

	return result
}

func (cm *ChannelManager) createCategory(categoryName string) (*discordgo.Channel, error) {
	category, err := cm.Session.GuildChannelCreateComplex(cm.GuildID, discordgo.GuildChannelCreateData{
		Name: categoryName,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	return category, err
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

// TODO this currently gets ALL channels. Maybe need a init method that puts into separate slices based on type (channel, category etc)
func getGuildChannelMap(s *discordgo.Session, guildID string) (map[channelKey]*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	channelMap := make(map[channelKey]*discordgo.Channel)
	for _, channel := range channels {
		channelMap[channelKey{Name: channel.Name, ParentID: channel.ParentID}] = channel
	}

	return channelMap, nil
}
