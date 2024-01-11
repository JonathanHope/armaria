package msgs

// ErrorMsg is a message that contains an error.
type ErrorMsg struct{ Err error }

// Error stringifies ErrorMsg.
func (e ErrorMsg) Error() string { return e.Err.Error() }

// View is which View to show.
type View string

const (
	ViewBooks View = "books" // show a listing of books
	ViewError View = "error" // show an error
)

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
