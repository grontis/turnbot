package game

type Stat struct {
	Name     string
	Modifier int
}

type Stats struct {
	Strength     Stat
	Dexterity    Stat
	Constitution Stat
	Intelligence Stat
	Wisdom       Stat
	Charisma     Stat
}
