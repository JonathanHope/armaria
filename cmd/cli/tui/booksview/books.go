package booksview

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/header"
	"github.com/jonathanhope/armaria/cmd/cli/tui/help"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/scrolltable"
	"github.com/jonathanhope/armaria/cmd/cli/tui/textinput"
	"github.com/jonathanhope/armaria/pkg/api"
	"github.com/jonathanhope/armaria/pkg/model"
)

const HeaderHeight = 3             // height of the header
const HeaderSpacerHeight = 1       // height of the spacer between the header and table
const FooterHeight = 3             // height of the footer
const HelpInfoWidth = 7            // width of the help info in the footer
const HeaderName = "BooksHeader"   // name of the header
const TextInputName = "BooksInput" // name of the textinput
const TableName = "BooksTable"     // name of the table

// InputType is which type of input is being collected.
type inputType int

const (
	inputNone   inputType = iota // not currently collecting input
	inputSearch                  // collecting input for a search
	inputURL                     // collecting input for an URL
	inputName                    // collecting input for a name
	inputFolder                  // collecting input to add a folder
)

// model is the model for the book listing.
// The book listing displays the bookmarks in the bookmarks DB.
type model struct {
	activeView msgs.View                                  // which view is currently being shown
	helpMode   bool                                       // whether to show the help or not
	inputMode  bool                                       // if true then the view is collecting input
	inputType  inputType                                  // which type of input is being collected
	width      int                                        // the current width of the screen
	height     int                                        // the current height of the screen
	folder     string                                     // the current folder
	query      string                                     // current search query
	busy       bool                                       // used to limit writers
	header     header.HeaderModel                         // header for app
	table      scrolltable.ScrolltableModel[armaria.Book] // table of books
	help       help.HelpModel                             // help for the app
	input      textinput.TextInputModel                   // allows text input
}

// InitialModel builds the model.
func InitialModel() tea.Model {
	return model{
		activeView: msgs.ViewBooks,
		header:     header.InitialModel(HeaderName, "📜 Armaria"),
		table: scrolltable.InitialModel[armaria.Book](
			TableName,
			[]scrolltable.ColumnDefinition[armaria.Book]{
				{
					Mode:        scrolltable.StaticColumn,
					StaticWidth: 4,
					RenderCell:  formatIsFolder,
					Style:       styleIsFolder,
				},
				{
					Mode:       scrolltable.DynamicColumn,
					Header:     "Name",
					RenderCell: formatName,
					Style:      styleURLNameTags,
				},
				{
					Mode:       scrolltable.DynamicColumn,
					Header:     "URL",
					RenderCell: formatURL,
					Style:      styleURLNameTags,
				},
				{
					Mode:       scrolltable.DynamicColumn,
					Header:     "Tags",
					RenderCell: formatTags,
					Style:      styleURLNameTags,
				},
			}),
		help: help.InitialModel(
			[]string{"Listing", "Input"},
			[]help.Binding{
				{Context: "Listing", Key: "up", Help: "Previous book"},
				{Context: "Listing", Key: "down", Help: "Next book"},
				{Context: "Listing", Key: "ctrl+up", Help: "Move book up"},
				{Context: "Listing", Key: "ctrl+down", Help: "Move book down"},
				{Context: "Listing", Key: "left", Help: "Move to parent folder"},
				{Context: "Listing", Key: "right", Help: "Move to folder children"},
				{Context: "Listing", Key: "enter", Help: "Open bookmark or folder"},
				{Context: "Listing", Key: "s", Help: "Search bookmarks/folders"},
				{Context: "Listing", Key: "c", Help: "Clear filters"},
				{Context: "Listing", Key: "r", Help: "Reload books"},
				{Context: "Listing", Key: "u", Help: "Edit URL"},
				{Context: "Listing", Key: "n", Help: "Edit name"},
				{Context: "Listing", Key: "+", Help: "Add folder"},
				{Context: "Listing", Key: "q", Help: "Quit"},
				{Context: "Input", Key: "left", Help: "Move to previous char"},
				{Context: "Input", Key: "right", Help: "Move to next char"},
				{Context: "Input", Key: "enter", Help: "Confirm input"},
				{Context: "Input", Key: "esc", Help: "Cancel input"},
			},
		),
		input: textinput.InitialModel(TextInputName, ""),
	}
}

// formatIsFolder formats an is folder value.
func formatIsFolder(book armaria.Book) string {
	if book.IsFolder {
		return "📁"
	}

	return "📖"
}

// styleIsFolder styles an is folder value.
func styleIsFolder(book armaria.Book, isSelected bool, isHeader bool) lipgloss.Style {
	return lipgloss.NewStyle().Align(lipgloss.Center)
}

// formatURL formats an URL value.
func formatURL(book armaria.Book) string {
	if book.URL != nil {
		return *book.URL
	}

	return ""
}

// styleURL styles an URL or Name value.
func styleURLNameTags(book armaria.Book, isSelected bool, isHeader bool) lipgloss.Style {
	if isHeader {
		return lipgloss.
			NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3"))
	}

	style := lipgloss.
		NewStyle()

	if book.IsFolder {
		style = style.Foreground(lipgloss.Color("4"))
	}

	if isSelected {
		style = style.Bold(true).Underline(true)
	}

	return style
}

// formatName formats a name value.
func formatName(book armaria.Book) string {
	return book.Name
}

// formatTags formats a tags value.
func formatTags(book armaria.Book) string {
	return strings.Join(book.Tags, ", ")
}

// filtersDisplayHeight returns the height of the filters display in the header.
func (m model) filtersDisplayHeight() int {
	if len(m.query) > 0 {
		return 1
	}

	return 0
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// There are some cases where this view is in a substate that handles events differently.

	// Substate: a different view is active.
	if m.activeView != msgs.ViewBooks {
		switch msg := msg.(type) {

		case msgs.ViewMsg:
			m.activeView = msgs.View(msg)
			return m, nil

		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
			return m, m.recalculateSizeCmd()

		default:
			return m, nil
		}
	}

	// Substate: the help screen is active.
	if m.helpMode {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit

			case "q", "esc":
				m.helpMode = false
				return m, nil

			default:
				return m, nil
			}

		case msgs.ViewMsg:
			m.activeView = msgs.View(msg)
			return m, nil

		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
			return m, m.recalculateSizeCmd()

		default:
			return m, nil
		}
	}

	// Substate: the input is active.
	if m.inputMode {
		switch msg := msg.(type) {

		case tea.KeyMsg:
			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit

			case "esc":
				m.inputMode = false
				m.query = ""
				m.inputType = inputNone
				return m, tea.Batch(m.inputEndCmd(), m.updateFiltersCmd())

			case "enter":
				cmds := []tea.Cmd{m.inputEndCmd()}
				if m.inputType == inputName {
					m.busy = true
					cmds = append(cmds, m.updateNameCmd(m.input.Text()))
				} else if m.inputType == inputURL {
					m.busy = true
					cmds = append(cmds, m.updateURLCmd(m.input.Text()))
				} else if m.inputType == inputFolder {
					m.busy = true
					cmds = append(cmds, m.addFolderCmd(m.input.Text()))
				}

				m.inputMode = false
				m.inputType = inputNone

				return m, tea.Batch(cmds...)

			default:
				var inputCmd tea.Cmd
				m.input, inputCmd = m.input.Update(msg)
				return m, inputCmd
			}

		case tea.WindowSizeMsg:
			m.height = msg.Height
			m.width = msg.Width
			return m, m.recalculateSizeCmd()

		case msgs.InputChangedMsg:
			if m.inputType == inputSearch {
				m.query = m.input.Text()
				return m, m.getBooksCmd(msgs.DirectionStart)
			}

		case msgs.DataMsg[armaria.Book]:
			var tableCmd tea.Cmd
			m.table, tableCmd = m.table.Update(msg)
			return m, tableCmd

		case msgs.ViewMsg:
			m.activeView = msgs.View(msg)
			return m, nil

		default:
			var inputCmd tea.Cmd
			m.input, inputCmd = m.input.Update(msg)
			return m, inputCmd
		}
	}

	//  Otherwise we fall into the main event loop.
	// The first step is forward the message to the underlying components.

	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)

	var headerCmd tea.Cmd
	m.header, headerCmd = m.header.Update(msg)

	cmds := []tea.Cmd{tableCmd, headerCmd, inputCmd}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			m.helpMode = true

		case "enter":
			if !m.table.Empty() {
				if m.table.Selection().IsFolder {
					m.folder = m.table.Selection().ID
					cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))
				} else {
					cmds = append(cmds, m.openURLCmd())
				}
			}

		case "left":
			if m.folder != "" {
				cmds = append(cmds, m.getParentCmd())
			}

		case "right":
			if !m.table.Empty() && m.table.Selection().IsFolder {
				m.folder = m.table.Selection().ID
				cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))
			}

		case "D":
			if !m.table.Empty() {
				m.busy = true
				cmds = append(cmds, m.deleteBookCmd())
			}

		case "c":
			m.query = ""
			cmds = append(
				cmds,
				m.getBooksCmd(msgs.DirectionStart),
				m.updateFiltersCmd(),
				m.recalculateSizeCmd(),
			)

		case "r":
			cmds = append(cmds, m.getBooksCmd(msgs.DirectionNone))

		case "s":
			m.inputMode = true
			m.inputType = inputSearch
			cmds = append(cmds, m.inputStartCmd("Query: ", ""))

		case "u":
			if !m.table.Empty() && !m.table.Selection().IsFolder {
				m.inputMode = true
				m.inputType = inputURL
				cmds = append(cmds, m.inputStartCmd("URL: ", *m.table.Selection().URL))
			}

		case "n":
			if !m.table.Empty() {
				m.inputMode = true
				m.inputType = inputName
				cmds = append(cmds, m.inputStartCmd("Name: ", m.table.Selection().Name))
			}

		case "+":
			m.inputMode = true
			m.inputType = inputFolder
			cmds = append(cmds, m.inputStartCmd("Folder: ", ""))

		case "ctrl+up":
			if m.query == "" && !m.table.Empty() && m.table.Index() > 0 && !m.busy {
				m.busy = true

				if m.table.Index() == 1 {
					next := m.table.Data()[0].ID
					cmds = append(cmds, m.moveToStartCmd(next, msgs.DirectionUp))
				} else {
					previous := m.table.Data()[m.table.Index()-2].ID
					next := m.table.Data()[m.table.Index()-1].ID
					cmds = append(cmds, m.moveBetweenCmd(previous, next, msgs.DirectionUp))
				}
			}

		case "ctrl+down":
			if m.query == "" && !m.table.Empty() && m.table.Index() < len(m.table.Data())-1 && !m.busy {
				m.busy = true

				if m.table.Index() == len(m.table.Data())-2 {
					previous := m.table.Data()[len(m.table.Data())-1].ID
					cmds = append(cmds, m.moveToEndCmd(previous, msgs.DirectionDown))
				} else {
					previous := m.table.Data()[m.table.Index()+1].ID
					next := m.table.Data()[m.table.Index()+2].ID
					cmds = append(cmds, m.moveBetweenCmd(previous, next, msgs.DirectionDown))
				}
			}
		}

	case msgs.SelectionChangedMsg[armaria.Book]:
		if !m.table.Empty() {
			cmds = append(cmds, m.getBreadcrumbsCmd())
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width

		cmds = append(cmds, m.recalculateSizeCmd())

	case msgs.FolderMsg:
		m.folder = string(msg)
		cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))

	case msgs.ViewMsg:
		m.activeView = msgs.View(msg)
		return m, nil

	case msgs.FreeMsg:
		m.busy = false
	}

	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m model) View() string {
	if m.activeView != msgs.ViewBooks {
		return ""
	}

	if m.helpMode {
		return m.header.View() + "\n\n" + m.help.View()
	}

	header := m.header.View()
	_, headerHeight := lipgloss.Size(header)

	table := m.table.View()
	_, tableHeight := lipgloss.Size(table)

	spacer := "\n\n"
	space := m.height -
		headerHeight -
		HeaderSpacerHeight -
		tableHeight -
		FooterHeight

	for i := 0; i < space; i++ {
		spacer += "\n"
	}

	return header +
		"\n\n" +
		table +
		spacer +
		m.footerView()
}

// footerView renders the footer view.
func (m model) footerView() string {
	help := lipgloss.
		NewStyle().
		Foreground(lipgloss.Color("3")).
		SetString("Help: ?")

	footerStyle := lipgloss.
		NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderTop(true).
		BorderBottom(false).
		BorderLeft(false).
		BorderRight(false).
		Width(m.width).
		BorderForeground(lipgloss.Color("5"))

	return footerStyle.Render(m.input.View() + help.Render())
}

// Init initializes the model.
func (m model) Init() tea.Cmd {
	return m.getBooksCmd(msgs.DirectionStart)
}

// getBooksCmd is a command to get books from the bookmarks database.
func (m model) getBooksCmd(move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.
			DefaultListBooksOptions().
			WithoutParentID()

		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		// The query must be at least 3 chars.
		if len(m.query) > 2 {
			options.WithQuery(m.query)
		}

		books, err := armariaapi.ListBooks(options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return msgs.DataMsg[armaria.Book]{Name: TableName, Data: books, Move: move}
	}
}

// getParentCmd is a command to go one level up in the folder structure.
func (m model) getParentCmd() tea.Cmd {
	return func() tea.Msg {
		book, err := armariaapi.GetBook(m.folder, armariaapi.DefaultGetBookOptions())
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		if book.ParentID == nil {
			return msgs.FolderMsg("")
		}

		return msgs.FolderMsg(*book.ParentID)
	}
}

// openURLCmd opens a bookmarks URL in the browser.
func (m model) openURLCmd() tea.Cmd {
	return func() tea.Msg {
		var cmd string
		var args []string

		switch runtime.GOOS {
		case "windows":
			cmd = "cmd"
			args = []string{"/c", "start"}
		case "darwin":
			cmd = "open"
		default: // "linux", "freebsd", "openbsd", "netbsd"
			cmd = "xdg-open"
		}

		args = append(args, *m.table.Selection().URL)
		err := exec.Command(cmd, args...).Start()
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return nil
	}
}

// getBreadcrumbsCmd gets breadcrumbs for the currently selected book.
func (m model) getBreadcrumbsCmd() tea.Cmd {
	return func() tea.Msg {
		parents, err := armariaapi.GetParentNames(m.table.Selection().ID, armariaapi.DefaultGetParentNameOptions())
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return msgs.NavMsg(strings.Join(parents, " > "))
	}
}

// deleteBookCmd deletes a bookmark or folder.
func (m model) deleteBookCmd() tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			err := armariaapi.RemoveFolder(m.table.Selection().ID, armariaapi.DefaultRemoveFolderOptions())
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			err := armariaapi.RemoveBook(m.table.Selection().ID, armariaapi.DefaultRemoveBookOptions())
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// updateURLCmd updates an URL for a bookmark.
func (m model) updateURLCmd(URL string) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.
			DefaultUpdateBookOptions().
			WithURL(URL)

		_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// updateNameCmd updates a name for a bookmark or folder.
func (m model) updateNameCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armariaapi.
				DefaultUpdateFolderOptions().
				WithName(name)

			_, err := armariaapi.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armariaapi.
				DefaultUpdateBookOptions().
				WithName(name)

			_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// addFolderCmd adds a folder to the bookmarks database.
func (m model) addFolderCmd(name string) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.DefaultAddFolderOptions()
		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		_, err := armariaapi.AddFolder(name, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// updateFiltersCmd will upate the filters display in the header based on the current filters.
func (m model) updateFiltersCmd() tea.Cmd {
	if len(m.query) > 0 {
		return func() tea.Msg { return msgs.FiltersMsg([]string{fmt.Sprintf("Query: %s", m.query)}) }
	}

	return func() tea.Msg { return msgs.FiltersMsg([]string{}) }
}

// inputStartCmd makes the necessary state updates when the input mode starts.
func (m model) inputStartCmd(prompt string, text string) tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return msgs.PromptMsg{Name: TextInputName, Prompt: prompt} },
		func() tea.Msg { return msgs.TextMsg{Name: TextInputName, Text: text} },
		func() tea.Msg { return msgs.FocusMsg{Name: TextInputName} },
	)
}

// inputEndCmd makes the necessary state updates when the input mode ends.
func (m model) inputEndCmd() tea.Cmd {
	cmds := []tea.Cmd{
		func() tea.Msg { return msgs.BlurMsg{Name: TextInputName} },
		func() tea.Msg { return msgs.PromptMsg{Name: TextInputName, Prompt: ""} },
		func() tea.Msg { return msgs.TextMsg{Name: TextInputName, Text: ""} },
		m.getBooksCmd(msgs.DirectionNone),
		m.updateFiltersCmd(),
		m.recalculateSizeCmd(),
	}

	return tea.Batch(cmds...)
}

// recalculateSizeCmd recalculates the size of the components.
// This needs to happen when the filters or the window size changes.
func (m model) recalculateSizeCmd() tea.Cmd {
	height := m.height -
		HeaderHeight -
		m.filtersDisplayHeight() -
		HeaderSpacerHeight -
		FooterHeight

	headerSizeMsg := msgs.SizeMsg{
		Name:  HeaderName,
		Width: m.width,
	}

	tableSizeMsg := msgs.SizeMsg{
		Name:   TableName,
		Width:  m.width,
		Height: height,
	}

	inputSizeMsg := msgs.SizeMsg{
		Name:  TextInputName,
		Width: m.width - HelpInfoWidth,
	}

	return tea.Batch(
		func() tea.Msg { return headerSizeMsg },
		func() tea.Msg { return tableSizeMsg },
		func() tea.Msg { return inputSizeMsg },
	)
}

// moveToEndCmd moves a bookmark or folder to the end of the list.
func (m model) moveToEndCmd(previous string, move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armariaapi.
				DefaultUpdateFolderOptions().
				WithOrderAfter(previous)

			_, err := armariaapi.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armariaapi.
				DefaultUpdateBookOptions().
				WithOrderAfter(previous)

			_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(move)()
	}
}

// moveToStartCmd moves a bookmark or folder to the end of the list.
func (m model) moveToStartCmd(next string, move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armariaapi.
				DefaultUpdateFolderOptions().
				WithOrderBefore(next)

			_, err := armariaapi.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armariaapi.
				DefaultUpdateBookOptions().
				WithOrderBefore(next)

			_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(move)()
	}
}

// moveBetweenCmd moves a bookmark or folder between two items on the list.
func (m model) moveBetweenCmd(previous string, next string, move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armariaapi.
				DefaultUpdateFolderOptions().
				WithOrderAfter(previous).
				WithOrderBefore(next)

			_, err := armariaapi.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				fmt.Println(err)
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armariaapi.
				DefaultUpdateBookOptions().
				WithOrderAfter(previous).
				WithOrderBefore(next)

			_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				fmt.Println(err)
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(move)()
	}
}
