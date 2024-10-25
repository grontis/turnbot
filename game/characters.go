package game

type Character struct {
	Name      string
	Age       int
	Level     int
	Race      Race
	Class     Class
	Stats     Stats
	Inventory Inventory
	Skills    []Skill
	Spells    []Spell
}

func NewCharacter(name string, age int, level int, class Class) *Character {
	return &Character{
		Name:      name,
		Age:       age,
		Level:     level,
		Class:     class,
		Stats:     Stats{},
		Inventory: Inventory{},
		Skills:    make([]Skill, 0),
		Spells:    make([]Spell, 0),
	}
}
