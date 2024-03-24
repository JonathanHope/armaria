package booksview

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/tui/footer"
	"github.com/jonathanhope/armaria/cmd/cli/tui/header"
	"github.com/jonathanhope/armaria/cmd/cli/tui/help"
	"github.com/jonathanhope/armaria/cmd/cli/tui/msgs"
	"github.com/jonathanhope/armaria/cmd/cli/tui/scrolltable"
	"github.com/jonathanhope/armaria/cmd/cli/tui/typeahead"
	"github.com/jonathanhope/armaria/pkg/api"
	"github.com/jonathanhope/armaria/pkg/model"
	"github.com/samber/lo"
)

const HeaderHeight = 3                 // height of the header
const HeaderSpacerHeight = 1           // height of the spacer between the header and table
const FooterHeight = 4                 // height of the footer
const HeaderName = "BooksHeader"       // name of the header
const FooterName = "BooksFooter"       // name of the footer
const TableName = "BooksTable"         // name of the table
const HelpName = "BooksHelp"           // name of the help screen
const TypeaheadName = "BooksTypeahead" // name of the typeahead
const AddTagOperation = "AddTag"       // operation to add a tag
const RemoveTagOperation = "RemoveTag" // operation to remove tag

// InputType is which type of input is being collected.
type inputType int

const (
	inputNone     inputType = iota // not currently collecting input
	inputSearch                    // collecting input for a search
	inputURL                       // collecting input for an URL
	inputName                      // collecting input for a name
	inputFolder                    // collecting input to add a folder
	inputBookmark                  // collecting input to add a bookmark
)

// model is the model for the book listing.
// The book listing displays the bookmarks in the bookmarks DB.
type model struct {
	activeView msgs.View                                  // which view is currently being shown
	inputType  inputType                                  // which type of input is being collected
	width      int                                        // the current width of the screen
	height     int                                        // the current height of the screen
	folder     string                                     // the current folder
	query      string                                     // current search query
	header     header.HeaderModel                         // header for app
	footer     footer.FooterModel                         // footer for app
	table      scrolltable.ScrolltableModel[armaria.Book] // table of books
	help       help.HelpModel                             // help for the app
	typeahead  typeahead.TypeaheadModel                   // typeahead for the app
}

// InitialModel builds the model.
func InitialModel() tea.Model {
	return model{
		activeView: msgs.ViewBooks,
		header:     header.InitialModel(HeaderName, "📜 Armaria"),
		footer:     footer.InitialModel(FooterName),
		table: scrolltable.InitialModel[armaria.Book](
			TableName,
			false,
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
			HelpName,
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
				{Context: "Listing", Key: "b", Help: "Add bookmark"},
				{Context: "Listing", Key: "t", Help: "Add tag"},
				{Context: "Listing", Key: "T", Help: "Remove tag"},
				{Context: "Listing", Key: "q", Help: "Quit"},
				{Context: "Input", Key: "left", Help: "Move to previous char"},
				{Context: "Input", Key: "right", Help: "Move to next char"},
				{Context: "Input", Key: "enter", Help: "Confirm input"},
				{Context: "Input", Key: "esc", Help: "Cancel input"},
			},
		),
		typeahead: typeahead.InitialModel(TypeaheadName),
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

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If another view is active ignore all keypresses.
	if _, ok := msg.(tea.KeyMsg); ok && m.activeView != msgs.ViewBooks {
		return m, nil
	}

	// If the help screen is active direct all keypresses to it.
	if _, ok := msg.(tea.KeyMsg); ok && m.help.HelpMode() {
		var helpCmd tea.Cmd
		m.help, helpCmd = m.help.Update(msg)
		return m, helpCmd
	}

	// If the footer is in input mode direct all keypresses to it.
	if _, ok := msg.(tea.KeyMsg); ok && m.footer.InputMode() {
		var footerCmd tea.Cmd
		m.footer, footerCmd = m.footer.Update(msg)
		return m, footerCmd
	}

	// If the typeahead is in input mode direct all keypresses to it.
	if _, ok := msg.(tea.KeyMsg); ok && m.typeahead.TypeaheadMode() {
		var typeaheadCmd tea.Cmd
		m.typeahead, typeaheadCmd = m.typeahead.Update(msg)
		return m, typeaheadCmd
	}

	var footerCmd tea.Cmd
	m.footer, footerCmd = m.footer.Update(msg)

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)

	var headerCmd tea.Cmd
	m.header, headerCmd = m.header.Update(msg)

	var helpCmd tea.Cmd
	m.help, helpCmd = m.help.Update(msg)

	var typeaheadCmd tea.Cmd
	m.typeahead, typeaheadCmd = m.typeahead.Update(msg)

	cmds := []tea.Cmd{tableCmd, headerCmd, footerCmd, helpCmd, typeaheadCmd}

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "?":
			cmds = append(cmds, func() tea.Msg { return msgs.ShowHelpMsg{Name: HelpName} })

		case "enter":
			if !m.table.Empty() && !m.header.Busy() {
				if m.table.Selection().IsFolder {
					m.folder = m.table.Selection().ID
					cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))
				} else {
					cmds = append(cmds, m.openURLCmd())
				}
			}

		case "left":
			if m.folder != "" && !m.header.Busy() {
				cmds = append(cmds, m.getParentCmd())
			}

		case "right":
			if !m.table.Empty() && m.table.Selection().IsFolder && !m.header.Busy() {
				m.folder = m.table.Selection().ID
				cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))
			}

		case "D":
			if !m.table.Empty() && !m.header.Busy() {
				cmds = append(cmds, m.deleteBookCmd(), func() tea.Msg { return msgs.BusyMsg{} })
			}

		case "c":
			if m.query != "" && !m.header.Busy() {
				m.query = ""
				cmds = append(
					cmds,
					m.getBooksCmd(msgs.DirectionStart),
					m.updateFiltersCmd(),
					m.recalculateSizeCmd(),
				)
			}

		case "r":
			if !m.header.Busy() {
				cmds = append(cmds, m.getBooksCmd(msgs.DirectionNone))
			}

		case "s":
			if !m.header.Busy() {
				m.inputType = inputSearch
				cmds = append(cmds, m.inputStartCmd("Query: ", "", 0))
			}

		case "u":
			if !m.table.Empty() && !m.table.Selection().IsFolder && !m.header.Busy() {
				m.inputType = inputURL
				cmds = append(cmds,
					m.inputStartCmd("URL: ", *m.table.Selection().URL, 2048),
					func() tea.Msg { return msgs.BusyMsg{} })
			}

		case "n":
			if !m.table.Empty() && !m.header.Busy() {
				m.inputType = inputName
				cmds = append(cmds,
					m.inputStartCmd("Name: ", m.table.Selection().Name, 2048),
					func() tea.Msg { return msgs.BusyMsg{} })
			}

		case "+":
			if !m.header.Busy() {
				m.inputType = inputFolder
				cmds = append(cmds,
					m.inputStartCmd("Folder: ", "", 2048),
					func() tea.Msg { return msgs.BusyMsg{} })
			}

		case "b":
			if !m.header.Busy() {
				m.inputType = inputBookmark
				cmds = append(cmds,
					m.inputStartCmd("Bookmark: ", "", 2048),
					func() tea.Msg { return msgs.BusyMsg{} })
			}

		case "ctrl+up":
			if m.query == "" && !m.table.Empty() && m.table.Index() > 0 && !m.header.Busy() {
				if m.table.Index() == 1 {
					next := m.table.Data()[0].ID
					cmds = append(cmds,
						m.moveToStartCmd(next, msgs.DirectionUp),
						func() tea.Msg { return msgs.BusyMsg{} })
				} else {
					previous := m.table.Data()[m.table.Index()-2].ID
					next := m.table.Data()[m.table.Index()-1].ID
					cmds = append(cmds,
						m.moveBetweenCmd(previous, next, msgs.DirectionUp),
						func() tea.Msg { return msgs.BusyMsg{} })
				}
			}

		case "ctrl+down":
			if m.query == "" &&
				!m.table.Empty() &&
				m.table.Index() < len(m.table.Data())-1 &&
				!m.header.Busy() {
				if m.table.Index() == len(m.table.Data())-2 {
					previous := m.table.Data()[len(m.table.Data())-1].ID
					cmds = append(cmds,
						m.moveToEndCmd(previous, msgs.DirectionDown),
						func() tea.Msg { return msgs.BusyMsg{} })
				} else {
					previous := m.table.Data()[m.table.Index()+1].ID
					next := m.table.Data()[m.table.Index()+2].ID
					cmds = append(cmds,
						m.moveBetweenCmd(previous, next, msgs.DirectionDown),
						func() tea.Msg { return msgs.BusyMsg{} })
				}
			}

		case "t":
			if !m.header.Busy() && !m.table.Empty() && !m.table.Selection().IsFolder {
				cmds = append(cmds, m.typeaheadStartCmd(
					"Add Tag: ",
					"",
					128,
					AddTagOperation,
					true,
					func() ([]string, error) {
						options := armariaapi.DefaultListTagsOptions()
						return armariaapi.ListTags(options)
					},
					func(query string) ([]string, error) {
						options := armariaapi.DefaultListTagsOptions().WithQuery(query)
						return armariaapi.ListTags(options)
					},
				))
			}

		case "T":
			cmds = append(cmds, m.typeaheadStartCmd(
				"Remove Tag: ",
				"",
				128,
				RemoveTagOperation,
				false,
				func() ([]string, error) {
					return m.table.Selection().Tags, nil
				},
				func(query string) ([]string, error) {
					tags := lo.Filter(m.table.Selection().Tags, func(tag string, index int) bool {
						return strings.Contains(tag, query)
					})

					return tags, nil
				},
			))
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

	case msgs.InputCancelledMsg:
		if msg.Name == FooterName {
			m.query = ""
			m.inputType = inputNone
			cmds = append(cmds, m.inputEndCmd())
		}

	case msgs.InputConfirmedMsg:
		if msg.Name == FooterName {
			cmds = append(cmds, m.inputEndCmd())
			if m.inputType == inputName {
				cmds = append(cmds, m.
					updateNameCmd(m.footer.Text()),
					func() tea.Msg { return msgs.BusyMsg{} })
			} else if m.inputType == inputURL {
				cmds = append(cmds,
					m.updateURLCmd(m.footer.Text()),
					func() tea.Msg { return msgs.BusyMsg{} })
			} else if m.inputType == inputFolder {
				cmds = append(cmds,
					m.addFolderCmd(m.footer.Text()),
					func() tea.Msg { return msgs.BusyMsg{} })
			} else if m.inputType == inputBookmark {
				cmds = append(cmds,
					m.addBookmarkCmd(m.footer.Text()),
					func() tea.Msg { return msgs.BusyMsg{} })
			}

			m.inputType = inputNone
		}

	case msgs.TypeaheadCancelledMsg:
		if msg.Name == TypeaheadName {
			cmds = append(cmds, m.typeaheadEndCmd())
		}

	case msgs.TypeaheadConfirmedMsg:
		if msg.Name == TypeaheadName && msg.Operation == AddTagOperation {
			cmds = append(cmds, m.typeaheadEndCmd(), m.addTag(msg.Value))
		} else if msg.Name == TypeaheadName && msg.Operation == RemoveTagOperation {
			cmds = append(cmds, m.typeaheadEndCmd(), m.removeTag(msg.Value))
		}

	case msgs.InputChangedMsg:
		if m.inputType == inputSearch {
			m.query = m.footer.Text()
			cmds = append(cmds, m.getBooksCmd(msgs.DirectionStart))
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the model.
func (m model) View() string {
	if m.activeView != msgs.ViewBooks {
		return ""
	}

	if m.help.HelpMode() {
		return m.header.View() + "\n\n" + m.help.View()
	}

	if m.typeahead.TypeaheadMode() {
		return m.header.View() + "\n\n" + m.typeahead.View()
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
		m.footer.View()
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
		options := armariaapi.DefaultGetParentNameOptions()
		parents, err := armariaapi.GetParentNames(m.table.Selection().ID, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return msgs.BreadcrumbsMsg(strings.Join(parents, " > "))
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

// addBookmarkCmd adds a bookmark to the bookmarks database.
func (m model) addBookmarkCmd(url string) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.DefaultAddBookOptions()
		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		_, err := armariaapi.AddBook(url, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// updateFiltersCmd will upate the filters display in the header based on the current filters.
func (m model) updateFiltersCmd() tea.Cmd {
	if len(m.query) > 0 {
		return func() tea.Msg {
			return msgs.FiltersMsg{Name: FooterName, Filters: []string{fmt.Sprintf("Query: %s", m.query)}}
		}
	}

	return func() tea.Msg {
		return msgs.FiltersMsg{Name: FooterName, Filters: []string{}}
	}
}

// inputStartCmd makes the necessary state updates when the input mode starts.
func (m model) inputStartCmd(prompt string, text string, maxChars int) tea.Cmd {
	return func() tea.Msg {
		return msgs.InputModeMsg{
			Name:      FooterName,
			InputMode: true,
			Prompt:    prompt,
			Text:      text,
			MaxChars:  maxChars,
		}
	}
}

// inputEndCmd makes the necessary state updates when the input mode ends.
func (m model) inputEndCmd() tea.Cmd {
	cmds := []tea.Cmd{
		func() tea.Msg {
			return msgs.InputModeMsg{
				Name:      FooterName,
				InputMode: false,
			}
		},
		m.getBooksCmd(msgs.DirectionNone),
		m.updateFiltersCmd(),
		m.recalculateSizeCmd(),
	}

	return tea.Batch(cmds...)
}

// typeaheadStartCmd makes the necessary state updates when the typeahead mode starts.
func (m model) typeaheadStartCmd(prompt string, text string, maxChars int, operation string, includeInput bool, unfilteredQuery func() ([]string, error), filteredQuery func(query string) ([]string, error)) tea.Cmd {
	return func() tea.Msg {
		return msgs.TypeaheadModeMsg{
			Name:            TypeaheadName,
			InputMode:       true,
			Prompt:          prompt,
			Text:            text,
			MaxChars:        maxChars,
			MinFilterChars:  3,
			Operation:       operation,
			IncludeInput:    includeInput,
			UnfilteredQuery: unfilteredQuery,
			FilteredQuery:   filteredQuery,
		}
	}
}

// typeaheadEndCmd makes the necessary state updates when the typeahead mode ends.
func (m model) typeaheadEndCmd() tea.Cmd {
	cmds := []tea.Cmd{
		func() tea.Msg {
			return msgs.TypeaheadModeMsg{
				Name:      TypeaheadName,
				InputMode: false,
			}
		},
		m.getBooksCmd(msgs.DirectionNone),
		m.recalculateSizeCmd(),
	}

	return tea.Batch(cmds...)
}

// recalculateSizeCmd recalculates the size of the components.
// This needs to happen when the filters or the window size changes.
func (m model) recalculateSizeCmd() tea.Cmd {
	tableHeight := m.height -
		HeaderHeight -
		HeaderSpacerHeight -
		FooterHeight

	typeaheadHeight := m.height -
		HeaderHeight -
		HeaderSpacerHeight

	headerSizeMsg := msgs.SizeMsg{
		Name:  HeaderName,
		Width: m.width,
	}

	footerSizeMsg := msgs.SizeMsg{
		Name:  FooterName,
		Width: m.width,
	}

	tableSizeMsg := msgs.SizeMsg{
		Name:   TableName,
		Width:  m.width,
		Height: tableHeight,
	}

	typeaheadSizeMsg := msgs.SizeMsg{
		Name:   TypeaheadName,
		Width:  m.width,
		Height: typeaheadHeight,
	}

	return tea.Batch(
		func() tea.Msg { return headerSizeMsg },
		func() tea.Msg { return footerSizeMsg },
		func() tea.Msg { return tableSizeMsg },
		func() tea.Msg { return typeaheadSizeMsg },
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
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armariaapi.
				DefaultUpdateBookOptions().
				WithOrderAfter(previous).
				WithOrderBefore(next)

			_, err := armariaapi.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(move)()
	}
}

// addTag adds a tag to a bookmark.
func (m model) addTag(tag string) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.DefaultAddTagsOptions()
		_, err := armariaapi.AddTags(m.table.Selection().ID, []string{tag}, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}
		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// removeTag removes a tag from a bookmark.
func (m model) removeTag(tag string) tea.Cmd {
	return func() tea.Msg {
		options := armariaapi.DefaultRemoveTagsOptions()
		_, err := armariaapi.RemoveTags(m.table.Selection().ID, []string{tag}, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}
		return m.getBooksCmd(msgs.DirectionNone)()
	}
}
