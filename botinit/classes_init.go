package botinit

//TODO rename botinit package to botrules or botlogic?

import "turnbot/game"

func Classes() []*game.Class {
	classes := make([]*game.Class, 0)

	classes = append(classes, game.NewClass("Wizard"))
	classes = append(classes, game.NewClass("Fighter"))
	classes = append(classes, game.NewClass("Rogue"))

	return classes
}
