package game

type Spell struct {
	Name        string
	Description string
	Handler     func()
}
