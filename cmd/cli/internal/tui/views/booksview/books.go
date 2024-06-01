package booksview

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/footer"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/header"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/help"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/scrolltable"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/controls/typeahead"
	"github.com/jonathanhope/armaria/cmd/cli/internal/tui/msgs"
	"github.com/jonathanhope/armaria/pkg"
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

// inputType is which type of input is being collected.
type inputType int

const (
	inputNone     inputType = iota // not currently collecting input
	inputSearch                    // collecting input for a search
	inputURL                       // collecting input for an URL
	inputName                      // collecting input for a name
	inputFolder                    // collecting input to add a folder
	inputBookmark                  // collecting input to add a bookmark
)

// operation is which typeahead operation is being executed.
type operation int

const (
	noOperation           operation = iota // not currently in a typeahead
	addTagOperation                        // using a typeahead to add a tag
	removeTagOperation                     // using a typeahead to remove a tag
	changeParentOperation                  // using a typeahead to change a books parent
)

// model is the model for the book listing.
// The book listing displays the bookmarks in the bookmarks DB.
type model struct {
	inputType inputType                                  // which type of input is being collected
	operation operation                                  // which typeahead operation is currently being executed
	width     int                                        // the current width of the screen
	height    int                                        // the current height of the screen
	folder    string                                     // the current folder
	query     string                                     // current search query
	header    header.HeaderModel                         // header for app
	footer    footer.FooterModel                         // footer for app
	table     scrolltable.ScrolltableModel[armaria.Book] // table of books
	help      help.HelpModel                             // help for the app
	typeahead typeahead.TypeaheadModel                   // typeahead for the app
}

// InitialModel builds the model.
func InitialModel() tea.Model {
	return model{
		header: header.InitialModel(HeaderName, "  Armaria"),
		footer: footer.InitialModel(FooterName),
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
				{Context: "Listing", Key: "p", Help: "Change parent"},
				{Context: "Listing", Key: "P", Help: "Remove parent"},
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

// resize changes the size of the books view.
func (m *model) resize() {
	tableHeight := m.height -
		HeaderHeight -
		HeaderSpacerHeight -
		FooterHeight

	typeaheadHeight := m.height -
		HeaderHeight -
		HeaderSpacerHeight

	m.header.Resize(m.width)
	m.footer.Resize(m.width)
	m.table.Resize(m.width, tableHeight)
	m.typeahead.Resize(m.width, typeaheadHeight)
}

// updateFilters will upate the filters display in the header based on the current filters.
func (m *model) updateFilters() {
	if len(m.query) > 0 {
		m.footer.SetFilters([]string{fmt.Sprintf("Query: %s", m.query)})
	} else {
		m.footer.SetFilters([]string{})
	}
}

// Update handles a message.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		m.resize()

	case msgs.DataMsg[armaria.Book]:
		if msg.Name == TableName {
			m.header.SetFree()
			return m, m.table.Reload(msg.Data, msg.Move)
		}

	case msgs.DataMsg[typeahead.TypeaheadItem]:
		if msg.Name == m.typeahead.TableName() {
			return m, m.typeahead.Reload(msg.Data, msg.Move)
		}

	case msgs.SelectionChangedMsg[armaria.Book]:
		if msg.Name == TableName && !m.table.Empty() {
			return m, m.getBreadcrumbsCmd()
		}

	case msgs.BreadcrumbsMsg:
		m.header.SetBreadcrumbs(string(msg))

	case msgs.FolderMsg:
		m.folder = string(msg)
		return m, m.getBooksCmd(msgs.DirectionStart)

	case msgs.InputChangedMsg:
		if msg.Name == m.footer.InputName() && m.inputType == inputSearch {
			m.query = m.footer.Text()
			return m, m.getBooksCmd(msgs.DirectionStart)
		} else if msg.Name == m.typeahead.InputName() {
			return m, m.typeahead.LoadItemsCmd()
		}

	case tea.KeyMsg:
		if m.footer.InputMode() {
			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit

			case "backspace":
				return m, m.footer.Delete()

			case "left":
				m.footer.MoveLeft()

			case "right":
				m.footer.MoveRight()

			case "esc":
				var cmd tea.Cmd
				m.query = ""
				m.footer.StopInputMode()

				if m.inputType == inputSearch {
					cmd = m.getBooksCmd(msgs.DirectionStart)
				}

				m.inputType = inputNone
				return m, cmd

			case "enter":
				var cmd tea.Cmd

				switch m.inputType {
				case inputName:
					m.header.SetBusy()
					cmd = m.updateNameCmd(m.footer.Text())
				case inputURL:
					m.header.SetBusy()
					cmd = m.updateURLCmd(m.footer.Text())
				case inputFolder:
					m.header.SetBusy()
					cmd = m.addFolderCmd(m.footer.Text())
				case inputBookmark:
					m.header.SetBusy()
					cmd = m.addBookmarkCmd(m.footer.Text())
				}

				m.footer.StopInputMode()
				m.updateFilters()
				m.inputType = inputNone

				return m, cmd

			default:
				return m, m.footer.Insert(msg.Runes)
			}
		} else if m.typeahead.TypeaheadMode() {
			switch msg.String() {

			case "ctrl+c":
				return m, tea.Quit

			case "left":
				m.typeahead.MoveLeft()

			case "right":
				m.typeahead.MoveRight()

			case "up":
				return m, m.typeahead.MoveUp()

			case "down":
				return m, m.typeahead.MoveDown()

			case "backspace":
				return m, m.typeahead.Delete()

			case "esc":
				m.header.SetFree()
				m.typeahead.StopTypeahead()

			case "enter":
				value := m.typeahead.Selection()
				m.typeahead.StopTypeahead()

				switch m.operation {
				case addTagOperation:
					return m, m.addTagCmd(value.Value)
				case removeTagOperation:
					return m, m.removeTagCmd(value.Value)
				case changeParentOperation:
					return m, m.changeParentCmd(value.Value)
				}

			default:
				return m, m.typeahead.Insert(msg.Runes)
			}

		} else if m.help.HelpMode() {
			switch msg.String() {

			case "q", "esc":
				m.help.HideHelp()

			case "ctrl+c":
				return m, tea.Quit
			}
		} else {
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "?":
				m.help.ShowHelp()

			case "up":
				return m, m.table.MoveUp()

			case "down":
				return m, m.table.MoveDown()

			case "left":
				if m.folder != "" && !m.header.Busy() {
					return m, m.getParentCmd()
				}

			case "right":
				if !m.table.Empty() && m.table.Selection().IsFolder && !m.header.Busy() {
					m.folder = m.table.Selection().ID
					return m, m.getBooksCmd(msgs.DirectionStart)
				}

			case "enter":
				if !m.table.Empty() && !m.header.Busy() {
					if m.table.Selection().IsFolder {
						m.folder = m.table.Selection().ID
						return m, m.getBooksCmd(msgs.DirectionStart)
					} else {
						return m, m.openURLCmd()
					}
				}

			case "r":
				if !m.header.Busy() {
					return m, m.getBooksCmd(msgs.DirectionNone)
				}

			case "D":
				if !m.table.Empty() && !m.header.Busy() {
					m.header.SetBusy()
					return m, m.deleteBookCmd()
				}

			case "s":
				if !m.header.Busy() {
					m.inputType = inputSearch
					m.footer.StartInputMode("Query: ", "", 0)
				}

			case "c":
				if m.query != "" && !m.header.Busy() {
					m.query = ""
					m.updateFilters()
					return m, m.getBooksCmd(msgs.DirectionNone)
				}

			case "u":
				if !m.table.Empty() && !m.table.Selection().IsFolder && !m.header.Busy() {
					m.inputType = inputURL
					m.footer.StartInputMode("URL: ", *m.table.Selection().URL, 2048)
					m.header.SetBusy()
				}

			case "n":
				if !m.table.Empty() && !m.header.Busy() {
					m.inputType = inputName
					m.footer.StartInputMode("Name: ", m.table.Selection().Name, 2048)
					m.header.SetBusy()
				}

			case "+":
				if !m.table.Empty() && !m.header.Busy() {
					m.inputType = inputFolder
					m.footer.StartInputMode("Folder: ", "", 2048)
					m.header.SetBusy()
				}

			case "b":
				if !m.table.Empty() && !m.header.Busy() {
					m.inputType = inputBookmark
					m.footer.StartInputMode("Bookmark: ", "", 2048)
					m.header.SetBusy()
				}

			case "ctrl+up":
				if m.query == "" && !m.table.Empty() && m.table.Index() > 0 && !m.header.Busy() {
					m.header.SetBusy()
					if m.table.Index() == 1 {
						next := m.table.Data()[0].ID
						return m, m.moveToStartCmd(next, msgs.DirectionUp)
					} else {
						previous := m.table.Data()[m.table.Index()-2].ID
						next := m.table.Data()[m.table.Index()-1].ID
						return m, m.moveBetweenCmd(previous, next, msgs.DirectionUp)
					}
				}

			case "ctrl+down":
				if m.query == "" && !m.table.Empty() && m.table.Index() < len(m.table.Data())-1 && !m.header.Busy() {
					m.header.SetBusy()
					if m.table.Index() == len(m.table.Data())-2 {
						previous := m.table.Data()[len(m.table.Data())-1].ID
						return m, m.moveToEndCmd(previous, msgs.DirectionDown)
					} else {
						previous := m.table.Data()[m.table.Index()+1].ID
						next := m.table.Data()[m.table.Index()+2].ID
						return m, m.moveBetweenCmd(previous, next, msgs.DirectionDown)
					}
				}

			case "t":
				if !m.header.Busy() && !m.table.Empty() && !m.table.Selection().IsFolder {
					m.operation = addTagOperation
					m.header.SetBusy()
					return m, m.typeahead.StartTypeahead(typeahead.StartTypeaheadPayload{
						Prompt:         "Add Tag: ",
						Text:           "",
						MaxChars:       128,
						IncludeInput:   true,
						MinFilterChars: 3,
						UnfilteredQuery: func() ([]typeahead.TypeaheadItem, error) {
							options := armaria.DefaultListTagsOptions()
							tags, err := armaria.ListTags(options)

							if err != nil {
								return nil, err
							}

							items := lo.Map(tags, func(tag string, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: tag, Value: tag}
							})

							return items, nil
						},
						FilteredQuery: func(query string) ([]typeahead.TypeaheadItem, error) {
							options := armaria.DefaultListTagsOptions().WithQuery(query)
							tags, err := armaria.ListTags(options)

							if err != nil {
								return nil, err
							}

							items := lo.Map(tags, func(tag string, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: tag, Value: tag}
							})

							return items, nil
						},
					})
				}

			case "T":
				if !m.header.Busy() && !m.table.Empty() && !m.table.Selection().IsFolder {
					m.operation = removeTagOperation
					m.header.SetBusy()
					return m, m.typeahead.StartTypeahead(typeahead.StartTypeaheadPayload{
						Prompt:         "Remove Tag: ",
						Text:           "",
						MaxChars:       128,
						IncludeInput:   false,
						MinFilterChars: 3,
						UnfilteredQuery: func() ([]typeahead.TypeaheadItem, error) {
							items := lo.Map(m.table.Selection().Tags, func(tag string, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: tag, Value: tag}
							})

							return items, nil
						},
						FilteredQuery: func(query string) ([]typeahead.TypeaheadItem, error) {
							tags := lo.Filter(m.table.Selection().Tags, func(tag string, index int) bool {
								return strings.Contains(tag, query)
							})

							items := lo.Map(tags, func(tag string, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: tag, Value: tag}
							})

							return items, nil
						},
					})
				}

			case "p":
				if !m.header.Busy() && !m.table.Empty() {
					m.operation = changeParentOperation
					m.header.SetFree()
					return m, m.typeahead.StartTypeahead(typeahead.StartTypeaheadPayload{
						Prompt:         "Change Parent: ",
						Text:           "",
						MaxChars:       2048,
						IncludeInput:   false,
						MinFilterChars: 3,
						UnfilteredQuery: func() ([]typeahead.TypeaheadItem, error) {
							options := armaria.DefaultListBooksOptions().WithFolders(true).WithBooks(false)
							books, err := armaria.ListBooks(options)

							if err != nil {
								return nil, err
							}

							items := lo.Map(books, func(book armaria.Book, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: book.Name, Value: book.ID}
							})

							return items, nil
						},
						FilteredQuery: func(query string) ([]typeahead.TypeaheadItem, error) {
							options := armaria.
								DefaultListBooksOptions().
								WithFolders(true).
								WithBooks(false).
								WithQuery(query)
							books, err := armaria.ListBooks(options)

							if err != nil {
								return nil, err
							}

							items := lo.Map(books, func(book armaria.Book, index int) typeahead.TypeaheadItem {
								return typeahead.TypeaheadItem{Label: book.Name, Value: book.ID}
							})

							return items, nil
						},
					})
				}

			case "P":
				if !m.header.Busy() && !m.table.Empty() && m.table.Selection().ParentID != nil {
					return m, m.removeParentCmd()
				}
			}
		}
	}

	return m, nil
}

// View renders the model.
func (m model) View() string {
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
func (m *model) getBooksCmd(move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		options := armaria.
			DefaultListBooksOptions().
			WithoutParentID()

		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		// The query must be at least 3 chars.
		if len(m.query) > 2 {
			options.WithQuery(m.query)
		}

		books, err := armaria.ListBooks(options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return msgs.DataMsg[armaria.Book]{Name: TableName, Data: books, Move: move}
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
		options := armaria.DefaultGetParentNameOptions()
		parents, err := armaria.GetParentNames(m.table.Selection().ID, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return msgs.BreadcrumbsMsg(strings.Join(parents, " > "))
	}
}

// getParentCmd is a command to go one level up in the folder structure.
func (m model) getParentCmd() tea.Cmd {
	return func() tea.Msg {
		book, err := armaria.GetBook(m.folder, armaria.DefaultGetBookOptions())
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		if book.ParentID == nil {
			return msgs.FolderMsg("")
		}

		return msgs.FolderMsg(*book.ParentID)
	}
}

// deleteBookCmd deletes a bookmark or folder.
func (m model) deleteBookCmd() tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			err := armaria.RemoveFolder(m.table.Selection().ID, armaria.DefaultRemoveFolderOptions())
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			err := armaria.RemoveBook(m.table.Selection().ID, armaria.DefaultRemoveBookOptions())
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
		options := armaria.
			DefaultUpdateBookOptions().
			WithURL(URL)

		_, err := armaria.UpdateBook(m.table.Selection().ID, options)
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
			options := armaria.
				DefaultUpdateFolderOptions().
				WithName(name)

			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.
				DefaultUpdateBookOptions().
				WithName(name)

			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
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
		options := armaria.DefaultAddFolderOptions()
		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		_, err := armaria.AddFolder(name, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// addBookmarkCmd adds a bookmark to the bookmarks database.
func (m model) addBookmarkCmd(url string) tea.Cmd {
	return func() tea.Msg {
		options := armaria.DefaultAddBookOptions()
		if m.folder != "" {
			options.WithParentID(m.folder)
		}

		_, err := armaria.AddBook(url, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// moveToEndCmd moves a bookmark or folder to the end of the list.
func (m model) moveToEndCmd(previous string, move msgs.Direction) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armaria.
				DefaultUpdateFolderOptions().
				WithOrderAfter(previous)

			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.
				DefaultUpdateBookOptions().
				WithOrderAfter(previous)

			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
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
			options := armaria.
				DefaultUpdateFolderOptions().
				WithOrderBefore(next)

			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.
				DefaultUpdateBookOptions().
				WithOrderBefore(next)

			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
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
			options := armaria.
				DefaultUpdateFolderOptions().
				WithOrderAfter(previous).
				WithOrderBefore(next)

			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.
				DefaultUpdateBookOptions().
				WithOrderAfter(previous).
				WithOrderBefore(next)

			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(move)()
	}
}

// removeParentCmd removes the parent of a bookmark or folder.
func (m model) removeParentCmd() tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armaria.DefaultUpdateFolderOptions().WithoutParentID()
			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.DefaultUpdateBookOptions().WithoutParentID()
			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// addTagCmd adds a tag to a bookmark.
func (m model) addTagCmd(tag string) tea.Cmd {
	return func() tea.Msg {
		options := armaria.DefaultAddTagsOptions()
		_, err := armaria.AddTags(m.table.Selection().ID, []string{tag}, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}
		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// removeTagCmd removes a tag from a bookmark.
func (m model) removeTagCmd(tag string) tea.Cmd {
	return func() tea.Msg {
		options := armaria.DefaultRemoveTagsOptions()
		_, err := armaria.RemoveTags(m.table.Selection().ID, []string{tag}, options)
		if err != nil {
			return msgs.ErrorMsg{Err: err}
		}
		return m.getBooksCmd(msgs.DirectionNone)()
	}
}

// changeParentCmd changes the parent of a bookmark or folder.
func (m model) changeParentCmd(parentID string) tea.Cmd {
	return func() tea.Msg {
		if m.table.Selection().IsFolder {
			options := armaria.DefaultUpdateFolderOptions().WithParentID(parentID)
			_, err := armaria.UpdateFolder(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		} else {
			options := armaria.DefaultUpdateBookOptions().WithParentID(parentID)
			_, err := armaria.UpdateBook(m.table.Selection().ID, options)
			if err != nil {
				return msgs.ErrorMsg{Err: err}
			}
		}

		return m.getBooksCmd(msgs.DirectionNone)()
	}
}
