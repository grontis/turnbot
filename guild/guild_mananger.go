package guild

import "github.com/bwmarrin/discordgo"

type GuildManager struct {
	Session        *discordgo.Session
	GuildID        string
	ChannelManager *ChannelManager
}

func NewGuildManager(s *discordgo.Session, guildID string) (*GuildManager, error) {
	channelManager, err := NewChannelManager(s, guildID)
	if err != nil {
		return nil, err
	}

	return &GuildManager{
		Session:        s,
		GuildID:        guildID,
		ChannelManager: channelManager,
	}, nil
}

func (gm *GuildManager) TryCreateCategory(categoryName string) (*discordgo.Channel, error) {
	category, err := gm.ChannelManager.TryCreateCategory(categoryName)
	return category, err
}

func (gm *GuildManager) TryCreateChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel, err := gm.ChannelManager.TryCreateChannelUnderCategory(channelName, categoryID)
	return channel, err
}
