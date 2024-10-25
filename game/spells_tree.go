package game

type SpellsTree struct {
	SpellsByLevel map[int]Spell
}

func NewSpellsTree() *SpellsTree {
	return &SpellsTree{
		SpellsByLevel: make(map[int]Spell),
	}
}
