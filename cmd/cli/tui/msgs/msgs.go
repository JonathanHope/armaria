package msgs

// Direction is a direction movement can occur in.
type Direction int

const (
	DirectionNone  Direction = iota // don't move
	DirectionUp                     // move up
	DirectionDown                   // move down
	DirectionStart                  // move to the start
)

// View is which View to show.
type View string

const (
	ViewBooks View = "books" // show a listing of books
	ViewError View = "error" // show an error
)

// ErrorMsg is a message that contains an error.
type ErrorMsg struct{ Err error }

// Error stringifies ErrorMsg.
func (e ErrorMsg) Error() string { return e.Err.Error() }

// ViewMsg is a message that contains which view to show.
type ViewMsg View

// SizeMsg is a message to inform the component it needs to resize.
type SizeMsg struct {
	Name   string // name of the target component
	Width  int    // the max width the component can occupy
	Height int    // the max height the component can occupy
}

// BusyMsg is a message that denotes the writer is busy.
type BusyMsg struct{}

// FreeMsg is a message that denotes the writer is free again.
type FreeMsg struct{}

// FolderMsg is a message that changes the current folder.
type FolderMsg string

// BreadcrumbsMsg is a message that changes the current breadcrumbs being displayed.
type BreadcrumbsMsg string

// FiltersMsg is published when the filters change.
type FiltersMsg struct {
	Name    string   // name of the target component
	Filters []string // current set of filters
}

// InputModeMsg is used to change to input mode.
type InputModeMsg struct {
	Name      string // name of the target component
	InputMode bool   // if true the footer will be in input mode
	Prompt    string // the prompt to show
	Text      string // the text to start the input with
	MaxChars  int    // the maximum number of chars to allow
}

// InputCancelledMsg is published when the user cancels input.
type InputCancelledMsg struct {
	Name string // name of the component which is collecting input.
}

// InputConfirmedMsg is published when the user confirms input.
type InputConfirmedMsg struct {
	Name string // name of the component which is collecting input.
}

// InputChangedMsg is published when the value in an input changes.
type InputChangedMsg struct {
	Name string // name of the target component
}

// FocusMsg is used to focus an input.
type FocusMsg struct {
	Name     string // name of the target component
	MaxChars int    // maximum number of chars to allow
}

// BlurMsg is published when the input blurs.
type BlurMsg struct {
	Name string // name of the target component
}

// TextMsg sets the text of an input.
type TextMsg struct {
	Name string // name of the target component
	Text string // text to set
}

// PromptMsg sets the prompt of an input.
type PromptMsg struct {
	Name   string // name of the target component
	Prompt string // prompt to use
}

// BlinkMsg is a message that causes the cursor to blink.
type BlinkMsg struct {
	Name string // name of the target component
}

// TypeaheadModeMsg is used to change to typeahead mode.
type TypeaheadModeMsg struct {
	Name            string                               // name of the target component
	InputMode       bool                                 // if true the typeahead will be in input mode
	Prompt          string                               // the prompt to show
	Text            string                               // the text to start the input with
	MaxChars        int                                  // the maximum number of chars to allow
	UnfilteredQuery func() ([]string, error)             // returns results when there isn't enough input
	FilteredQuery   func(query string) ([]string, error) // returns results when there's enough input
	MinFilterChars  int                                  // the minumum number of chars needed to filter
	Operation       string                               // the operation the typeahead is for
	IncludeInput    bool                                 // if true include the current input as an option
}

// TypeaheadConfirmedMsg is published when an option is selected..
type TypeaheadConfirmedMsg struct {
	Name      string // name of the target typeahead
	Value     string // the value that was selected
	Operation string // the operation the typeahead is for
}

// SelectionConfirmedMsg is published when an option selection is cancelled.
type TypeaheadCancelledMsg struct {
	Name string // name of the target typeahead
}

// ShowHelpMsg is used to bring up a help screen.
type ShowHelpMsg struct {
	Name string // name of the target help screen
}

// SelectionChangedMsg is published when the table selection changes.
type SelectionChangedMsg[T any] struct {
	Name      string // name of the target component
	Empty     bool   // whether the frame is empty
	Selection T      // the selected item (if any) in the frame
}

// DataMsg is a message to update the data in the table.
type DataMsg[T any] struct {
	Name string    // name of the target component
	Data []T       //the data to show in the scrolltable
	Move Direction // optionally adjust the cursor
}
