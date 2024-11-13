package game

import "fmt"

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

func (c *Character) ToMessageContent() string {
	message := ""
	if c.Name != "" {
		message += fmt.Sprintf("# %s\n", c.Name)
	}

	if c.Class.Name != "" {
		message += fmt.Sprintf("### %s\n", c.Class.Name)
	}

	if c.Race.Name != "" {
		message += fmt.Sprintf("### %s\n", c.Race.Name)
	}
	//TODO other class properties

	return message
}

type CharacterManager struct {
	PlayerCharacters map[string]*Character
}

func NewCharacterManager() *CharacterManager {
	return &CharacterManager{
		PlayerCharacters: make(map[string]*Character),
	}
}

func (cm *CharacterManager) AddNewCharacter(userID string, character *Character) {
	cm.PlayerCharacters[userID] = character
}

func (cm *CharacterManager) UpdateCharacterInfo(userID string, name string, age int) error {
	character := cm.PlayerCharacters[userID]
	if character == nil {
		return fmt.Errorf("no character found to update for userID: %s", userID)
	}

	character.Name = name
	character.Age = age
	return nil
}

func (cm *CharacterManager) UpdateCharacterClass(userID string, className string) error {
	character := cm.PlayerCharacters[userID]
	if character == nil {
		return fmt.Errorf("no character found to update for userID: %s", userID)
	}

	character.Class = *NewClass(className)
	return nil
}

func (cm *CharacterManager) UpdateCharacterRace(userID string, raceName string) error {
	character := cm.PlayerCharacters[userID]
	if character == nil {
		return fmt.Errorf("no character found to update for userID: %s", userID)
	}

	character.Race = Race{Name: raceName}
	return nil
}
