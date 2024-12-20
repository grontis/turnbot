package game

type InteractionsInitLoader interface {
	LoadButtonInteractions(engine *GameEngine)
	LoadCommandInteractions(engine *GameEngine)
	LoadDropdownInteractions(engine *GameEngine)
	LoadModalInteractions(engine *GameEngine)
	LoadInteractionsHandler(engine *GameEngine)
	CreateAllCommands(engine *GameEngine)
}

type GuildInitLoader interface {
	SetupBotChannels(engine *GameEngine, guildID string) error
}
