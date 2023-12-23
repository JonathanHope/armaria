package msgs

// SelectionChangedMsg is published when the scrolltable selection changes.
type SelectionChangedMsg[T any] struct {
	Name      string // name of the target srolltable
	Empty     bool   // whether the frame is empty
	Selection T      // the selected item (if any) in the frame
}

// DataMsg is a message to update the data in the scrolltable.
type DataMsg[T any] struct {
	Name string    // name of the target scrolltable
	Data []T       //the data to show in the scrolltable
	Move Direction // optionally adjust the cursor
}

// Direction is a direction the cursor can move on the scrolltable..
type Direction int

const (
	DirectionNone  Direction = iota // don't move
	DirectionUp                     // move up the table
	DirectionDown                   // move down the table
	DirectionStart                  // move to the start of the table
)
