package msgs

// FocusMsg is used to focus a textinput.
type FocusMsg struct {
	Name string // name of the target textinput
}

// InputChangedMsg is published when the value in a textinput changes.
type InputChangedMsg struct {
	Name string // name of the target textinput
}

// BlinkMsg is a message that causes the cursor to blink.
type BlinkMsg struct {
	Name string // name of the target textinput
}

// BlurMsg is published when the textinput blurs.
type BlurMsg struct {
	Name string // name of the target textinput
}

// TextMsg sets the text of a textinput.
type TextMsg struct {
	Name string // name of the target textinput
	Text string // text to set on the textinput
}

// PromptMsg sets the prompt of a textinput.
type PromptMsg struct {
	Name   string // name of the target textinput
	Prompt string // text to set on the textinput
}
