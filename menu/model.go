package menu

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevm/bubbleo/navstack"
	"github.com/kevm/bubbleo/styles"
)

type Choice struct {
	Title       string
	Description string
	Model       tea.Model
}

type choiceItem struct {
	title, desc string
	key         Choice
}

func (i choiceItem) Title() string       { return i.title }
func (i choiceItem) Description() string { return i.desc }
func (i choiceItem) FilterValue() string { return i.title + i.desc }

type Model struct {
	Choices []Choice
	list    list.Model

	selected *Choice
}

// New setups up a new menu model
func New(title string, choices []Choice, selected *Choice) Model {
	delegation := list.NewDefaultDelegate()
	items := make([]list.Item, len(choices))
	selectedIndex := -1
	for i, choice := range choices {
		if selected != nil && &choice == selected {
			selectedIndex = i
		}
		items[i] = choiceItem{title: choice.Title, desc: choice.Description, key: choice}
	}

	model := Model{
		Choices:  choices,
		list:     list.New(items, delegation, 120, 20),
		selected: selected,
	}

	if selected != nil {
		model.list.Select(selectedIndex)
	}

	model.list.Styles.Title = styles.ListTitleStyle
	model.list.Title = title
	model.list.SetShowPagination(true)
	model.list.SetShowTitle(true)
	model.list.SetFilteringEnabled(false)
	model.list.SetShowFilter(false)
	model.list.SetShowStatusBar(false)
	model.list.SetShowHelp(false)

	//TODO: figure out height long term.
	// model.list.SetSize(window.Width, window.Height-window.TopOffset)

	chooseKeyBinding := key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	)
	model.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{chooseKeyBinding}
	}
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{chooseKeyBinding}
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.SetSize(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyEsc.String():
			pop := cmdize(navstack.PopNavigation{})
			return m, pop
		case tea.KeyEnter.String():
			choice, ok := m.list.SelectedItem().(choiceItem)
			if ok {
				m.selected = &choice.key
				item := navstack.NavigationItem{Title: choice.title, Model: choice.key.Model}
				cmd := cmdize(navstack.PushNavigation{Item: item})
				return m, cmd
			}
		}
	}

	// No selection made yet so update the list
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) SetSize(w tea.WindowSizeMsg) {
	m.list.SetSize(w.Width, w.Height)
}

func (m Model) View() string {
	// display menu if choices are present.
	if len(m.Choices) > 0 {
		return "\n" + m.list.View()
	}

	return ""
}

func cmdize[T any](t T) tea.Cmd {
	return func() tea.Msg {
		return t
	}
}
