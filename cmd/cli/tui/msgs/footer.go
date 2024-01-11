package msgs

// FocusMsg is used to set whether the footer is in input mode or not.
type InputModeMsg struct {
	Name      string // name of the target footer
	InputMode bool   // if true the footer will be in input mode
	Prompt    string // the prompt to show
	Text      string // the text to start the input with
	MaxChars  int    // the maximum number of chars to allow
}

// InputCancelledMsg is published when the user cancels input.
type InputCancelledMsg struct{}

// InputConfirmedMsg is published when the user cancels input.
type InputConfirmedMsg struct{}

// FiltersMsg is published when the filters change.
type FiltersMsg struct {
	Name    string   // name of the target footer
	Filters []string // current set of filters
}
