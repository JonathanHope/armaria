package msgs

// SelectionChangedMsg is published when the scrolltable selection changes.
type SelectionChangedMsg[T any] struct {
	Name      string // name of the target srolltable
	Empty     bool   // whether the frame is empty
	Selection T      // the selected item (if any) in the frame
}

// DataMsg is a message to update the data in the scrolltable.
type DataMsg[T any] struct {
	Name string // name of the target scrolltable
	Data []T    //the data to show in the scrolltable
}
