package game

type Class struct {
	Name       string
	SkillsTree *SkillsTree
	SpellsTree *SpellsTree
}

func NewClass(name string) *Class {
	return &Class{
		Name:       name,
		SkillsTree: NewSkillsTree(),
		SpellsTree: NewSpellsTree(),
	}
}
