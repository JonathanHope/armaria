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
type View int

const (
	ViewNone  View = iota // no view
	ViewBooks             // show a listing of books
	ViewError             // show an error
)

// ErrorMsg is a message that contains an error.
type ErrorMsg struct{ Err error }

// Error stringifies ErrorMsg.
func (e ErrorMsg) Error() string { return e.Err.Error() }

// ViewMsg is a message that contains which view to show.
type ViewMsg View

// SelectionChangedMsg is published when the table selection changes.
type SelectionChangedMsg[T any] struct {
	Name      string // name of the target control
	Empty     bool   // whether the frame is empty
	Selection T      // the selected item (if any) in the frame
}

// DataMsg is a message to update the data in the table.
type DataMsg[T any] struct {
	Name string    // name of the target control
	Data []T       //the data to show in the scrolltable
	Move Direction // optionally adjust the cursor
}

// InputChangedMsg is published when the value in an input changes.
type InputChangedMsg struct {
	Name string // name of the target control
}

// BreadcrumbsMsg is a message that changes the current breadcrumbs being displayed.
type BreadcrumbsMsg string

// FolderMsg is a message that changes the current folder.
type FolderMsg string
