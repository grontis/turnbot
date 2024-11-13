package guild

import "github.com/bwmarrin/discordgo"

type GuildManager struct {
	Session        *discordgo.Session
	GuildID        string
	channelManager *channelManager
}

func NewGuildManager(s *discordgo.Session, guildID string) (*GuildManager, error) {
	channelManager, err := newChannelManager(s, guildID)
	if err != nil {
		return nil, err
	}

	return &GuildManager{
		Session:        s,
		GuildID:        guildID,
		channelManager: channelManager,
	}, nil
}

func (gm *GuildManager) UserByID(userID string) (*discordgo.User, error) {
	member, err := gm.Session.GuildMember(gm.GuildID, userID)
	if err != nil {
		return nil, err
	}
	return member.User, nil
}

func (gm *GuildManager) FindCategoryByName(categoryName string) (*discordgo.Channel, error) {
	category, err := gm.channelManager.findCategoryByName(categoryName)
	return category, err
}

func (gm *GuildManager) FindChannelInCategoryByName(categoryName string, channelName string) (*discordgo.Channel, error) {
	channel, err := gm.channelManager.findChannelInCategoryByName(categoryName, channelName)
	return channel, err
}

func (gm *GuildManager) TryCreateCategory(categoryName string) (*discordgo.Channel, error) {
	category, err := gm.channelManager.tryCreateCategory(categoryName)
	return category, err
}

func (gm *GuildManager) TryCreateChannelUnderCategory(channelName string, categoryID string) (*discordgo.Channel, error) {
	channel, err := gm.channelManager.tryCreateChannelUnderCategory(channelName, categoryID)
	return channel, err
}
