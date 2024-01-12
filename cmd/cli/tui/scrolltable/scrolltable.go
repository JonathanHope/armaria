package scrolltable

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/utils"
	"github.com/samber/lo"
)

const Reserved = 4 // how much space things like headers and borders take up

// ColumnMode configures the width behavior of a column in the scrolltable.
type ColumnMode int

const (
	StaticColumn  ColumnMode = iota // the columns width is always the same
	DynamicColumn                   // the column will take up as much space as it can
)

// StyleSelectFn is a function to get the correct style for a given cell.
type StyleSelectFn[T any] func(datum T, isSelected bool, isHeader bool) lipgloss.Style

// RenderCellFn is a function that renders the string content for a given cell.
type RenderCellFn[T any] func(datum T) string

// ColumnDefinition defines the behavior for a given column in the scrolltable.
type ColumnDefinition[T any] struct {
	Mode        ColumnMode       // the width behavior of the column
	StaticWidth int              // static width of the column; only valid for StaticColumn
	Header      string           // header text
	Style       StyleSelectFn[T] // function to get the style for a given cell
	RenderCell  RenderCellFn[T]  // function to get the string content for a given cell
}

// ScrolltableModel is the ScrolltableModel for a scrolltable.
// The scrolltable is a table that supports data larger than the screen itself.
type ScrolltableModel[T any] struct {
	name              string                // name of the scrolltable
	columnDefinitions []ColumnDefinition[T] // information about the columns needed to style/render them
	data              []T                   //  the data to show in the scrolltable
	width             int                   // max width the scrolltable can occupy
	height            int                   // max height the scrolltable can occupy
	cursor            int                   // index of selected datum in the visible frame
	frameStart        int                   // start of the visible frame
}

// InitialModel builds a scrolltable model.
func InitialModel[T any](name string, columnDefinitions []ColumnDefinition[T]) ScrolltableModel[T] {
	return ScrolltableModel[T]{
		name:              name,
		columnDefinitions: columnDefinitions,
	}
}

// Empty returns true if the current frame is Empty.
func (m ScrolltableModel[T]) Empty() bool {
	return len(m.Frame()) == 0
}

// Selection returns the current Selection.
func (m ScrolltableModel[T]) Selection() T {
	var zero T
	if m.Empty() {
		return zero
	}

	return m.Frame()[m.cursor]
}

// Index returns the absolute Index of the current selection.
func (m ScrolltableModel[T]) Index() int {
	return m.frameStart + m.cursor
}

// Frame returns the currently visible Frame of data.
func (m ScrolltableModel[T]) Frame() []T {
	if len(m.data) == 0 {
		return nil
	}

	return m.data[m.frameStart : m.frameStart+m.frameSize()]
}

// Data returns the entire data source.
func (m ScrolltableModel[T]) Data() []T {
	return m.data
}

// frameSize returns the current size of the visible frame of data.
func (m ScrolltableModel[T]) frameSize() int {
	frameSize := m.height - Reserved
	if len(m.data) < frameSize {
		frameSize = len(m.data)
	}

	return frameSize
}

// resetFrame resets the frame after the size or data has changed.
func (m *ScrolltableModel[T]) resetFrame(move msgs.Direction) {
	if move == msgs.DirectionUp {
		m.moveUp()
	} else if move == msgs.DirectionDown {
		m.moveDown()
	} else if move == msgs.DirectionStart {
		m.cursor = 0
		m.frameStart = 0
	}

	if m.frameStart+m.cursor >= len(m.data) {
		m.cursor = 0
		m.frameStart = 0
	}
}

// moveDown moves the cursor down the table.
func (m *ScrolltableModel[T]) moveDown() (ScrolltableModel[T], tea.Cmd) {
	if m.Empty() {
		return *m, nil
	}

	move := false
	scroll := false
	if m.cursor != m.frameSize()-1 {
		// If the cursor isn't already at the end of the frame move it down.
		m.cursor += 1
		move = true
	} else if len(m.data) != m.frameStart+m.frameSize() {
		// If cursor is at the end of the frame and not at the end of data move the frame down.
		scroll = true
		m.frameStart += 1
	}

	if scroll || move {
		return *m, m.selectionChangedCmd()
	}

	return *m, nil
}

// moveUp moves the cursor up the table.
func (m *ScrolltableModel[T]) moveUp() (ScrolltableModel[T], tea.Cmd) {
	if m.Empty() {
		return *m, nil
	}

	move := false
	scroll := false
	if m.cursor != 0 {
		// If the cursor isn't already at the start of the frame move it up.
		m.cursor -= 1
		move = true
	} else if m.frameStart != 0 {
		// If cursor is at the start of the frame and not at the start of the data move the frame up.
		scroll = true
		m.frameStart -= 1
	}

	if scroll || move {
		return *m, m.selectionChangedCmd()
	}

	return *m, nil
}

// Update updates the scrolltable model from a message.
func (m ScrolltableModel[T]) Update(msg tea.Msg) (ScrolltableModel[T], tea.Cmd) {
	switch msg := msg.(type) {

	case msgs.SizeMsg:
		if msg.Name == m.name {
			m.width = msg.Width
			m.height = msg.Height
			m.resetFrame(msgs.DirectionNone)
		}

	case msgs.DataMsg[T]:
		if msg.Name == m.name {
			// This allows the cursor to stick in the right place when an item is removed.
			if len(m.data)-1 == len(msg.Data) && m.frameStart > 0 {
				m.frameStart -= 1
			}

			m.data = msg.Data
			m.resetFrame(msg.Move)
			return m, tea.Batch(m.selectionChangedCmd(), func() tea.Msg { return msgs.FreeMsg{} })
		}

	case tea.KeyMsg:
		switch msg.String() {

		case "down":
			return m.moveDown()

		case "up":
			return m.moveUp()
		}
	}

	return m, nil
}

// View renders the scrolltable model.
func (m ScrolltableModel[T]) View() string {
	const cellPadding = 1

	if m.Empty() {
		return ""
	}

	// Columns can have static or dynamic width.
	// Dynamic columns occupy the width less the static widths.
	// Each dynamic column gets the same amount of space.

	staticWidths := 0
	numberOfDynamicColumns := 0
	for _, def := range m.columnDefinitions {
		if def.Mode == StaticColumn {
			staticWidths += def.StaticWidth
		} else {
			numberOfDynamicColumns += 1
		}
	}

	numberOfVerticalBorders := len(m.columnDefinitions) - 1
	availableWidth := m.width - (numberOfVerticalBorders + staticWidths)
	dynamicColumnWidth := availableWidth / numberOfDynamicColumns
	dynamicColumnTextWidth := dynamicColumnWidth - cellPadding*2

	rows := [][]string{}
	lo.ForEach(m.Frame(), func(datum T, _ int) {
		row := lo.Map(m.columnDefinitions, func(def ColumnDefinition[T], _ int) string {
			cell := def.RenderCell(datum)
			if def.Mode == DynamicColumn {
				// Dynamic columns need to have their string truncated if it's too long.
				cell = utils.Substr(cell, 0, dynamicColumnTextWidth)
			}
			return cell
		})

		rows = append(rows, row)
	})

	headers := lo.Map(m.columnDefinitions, func(def ColumnDefinition[T], _ int) string {
		return def.Header
	})

	booksTable := table.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(false).
		BorderRight(false).
		BorderColumn(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("6"))).
		Width(m.width).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			def := m.columnDefinitions[col]
			var datum T
			if row > 0 {
				datum = m.Frame()[row-1]
			}

			style := m.columnDefinitions[col].Style(datum, row == m.cursor+1, row == 0)
			if def.Mode == StaticColumn {
				return style.Width(def.StaticWidth).Padding(0, cellPadding)
			} else {
				return style.Width(dynamicColumnWidth).Padding(0, cellPadding)
			}
		}).
		Rows(rows...)

	return booksTable.Render()
}

// Init initializes the scrolltable model.
func (m ScrolltableModel[T]) Init() tea.Cmd {
	return nil
}

// selectionChangedCmd publishes a message with the currently selected datum.
func (m ScrolltableModel[T]) selectionChangedCmd() tea.Cmd {
	return func() tea.Msg {
		return msgs.SelectionChangedMsg[T]{
			Empty:     m.Empty(),
			Selection: m.Selection(),
		}
	}
}
