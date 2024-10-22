package game

type InteractionLogicLoader interface {
	LoadButtonInteractions(engine *GameEngine)
	LoadCommandInteractions(engine *GameEngine)
	LoadDropdownInteractions(engine *GameEngine)
	LoadModalInteractions(engine *GameEngine)
	LoadInteractionsHandler(engine *GameEngine)
}
