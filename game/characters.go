package game

type Character struct {
	Name      string
	Age       int
	Class     Class
	Stats     Stats
	Inventory Inventory
	Skills    []Skill
	Spells    []Spell
}

func NewCharacter(name string, age int, class Class) *Character {
	return &Character{
		Name:      name,
		Age:       age,
		Class:     class,
		Stats:     Stats{},
		Inventory: Inventory{},
		Skills:    make([]Skill, 0),
		Spells:    make([]Spell, 0),
	}
}
