package game

type SkillsTree struct {
	SkillsByLevel map[int]Skill
}

func NewSkillsTree() *SkillsTree {
	return &SkillsTree{
		SkillsByLevel: make(map[int]Skill),
	}
}
