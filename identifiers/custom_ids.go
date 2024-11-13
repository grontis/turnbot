package identifiers

type ButtonCustomID string

const (
	ButtonStartCharacterCreationCustomID ButtonCustomID = "button_start_character_creation"
	ButtonDiceRollCustomID               ButtonCustomID = "button_dice_roll"
	ButtonOpenCharacterInfoModalCustomID ButtonCustomID = "button_open_character_info_modal"
)

type CommandNameID string

const (
	CommandHello CommandNameID = "command_hello"
)

type DropdownCustomID string

const (
	DropdownClassSelectCustomID DropdownCustomID = "dropdown_class_select"
	DropdownRaceSelectCustomID  DropdownCustomID = "dropdown_race_select"
)

type ModalCustomID string

const (
	ModalCharacterInfoCustomID ModalCustomID = "model_character_info"
)

type TextInputCustomID string

const (
	TextInputCharacterName TextInputCustomID = "textinput_character_name"
	TextInputCharacterAge  TextInputCustomID = "textinput_character_age"
)
