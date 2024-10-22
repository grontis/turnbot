package guild

import "github.com/bwmarrin/discordgo"

type GuildManager struct {
	Session        *discordgo.Session
	GuildID        string
	ChannelManager *channelManager
}

func NewGuildManager(s *discordgo.Session, guildID string) (*GuildManager, error) {
	channelManager, err := newChannelManager(s, guildID)
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
	category, err := gm.ChannelManager.tryCreateCategory(categoryName)
	return category, err
}

func (gm *GuildManager) TryCreateChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel, err := gm.ChannelManager.tryCreateChannelUnderCategory(channelName, categoryID)
	return channel, err
}
