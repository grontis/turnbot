package game

type Item struct {
	Name           string
	Description    string
	Value          int
	AbilityHandler func()
}
