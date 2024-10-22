package botlogic

import (
	"fmt"
	"log"
	"turnbot/game"
)

type BotGuildLogicLoader struct{}

func (b *BotGuildLogicLoader) SetupBotChannels(engine *game.GameEngine, guildID string) error {
	turnbotCategoryName := "turnbot"
	turnbotCategory, err := engine.GuildManager.TryCreateCategory(turnbotCategoryName)
	if err != nil {
		log.Printf("error creating category %s: %s", turnbotCategoryName, err)
	}

	_, err = engine.GuildManager.TryCreateChannelUnderCategory("bot-test", turnbotCategory.ID)
	if err != nil {
		fmt.Printf("error creating channel: %s", err)
		return err
	}

	return nil
}
