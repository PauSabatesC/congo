package tui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/PauSabatesC/congo/congo"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14
const defaultWidth = 20

var Exited bool = false

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	nextTextStyle     = lipgloss.NewStyle().Margin(1, 0, 1, 2)
)

type item struct {
	title       string
	description string
}

func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type model struct {
	list     list.Model
	items    []item
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			Exited = true
			m.quitting = true
			return m, tea.Quit

		case "q":
			Exited = true
			m.quitting = true
			return m, tea.Quit

		case "esc":
			Exited = true
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.title)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return nextTextStyle.Render(fmt.Sprintf("-> Selected: %s", m.choice))
	}
	if m.quitting {
		return quitTextStyle.Render("Exited correctly")
	}
	return "\n" + m.list.View()
}

func PrintList(elements []string, title string) (string, error) {
	var items []list.Item
	for _, element := range elements {
		var itemAux item
		itemAux.title = element
		items = append(items, itemAux)
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	p := tea.NewProgram(model{list: l})

	defer tea.Quit()

	m, err := p.StartReturningModel()
	if err != nil {
		return "", err
	}

	var choice string
	m, ok := m.(model)
	if ok && m.(model).choice != "" {
		choice = m.(model).choice
	}

	return choice, nil
}

func SelectElementFromList(elements []string, title string) (string, error) {
	if len(elements) == 0 {
		return "", errors.New("Empty list")
	}

	choice, err := PrintList(elements, title)
	if err != nil {
		return "", err
	}

	return choice, nil
}

func SelectElementFromEcsTask(ecsTasks []congo.EcsTaskData, title string) (*congo.EcsTaskData, error) {
	if len(ecsTasks) == 0 {
		return nil, errors.New("Empty list")
	}
	var taskIds []string
	for _, task := range ecsTasks {
		taskIds = append(taskIds, task.Id+" | "+task.StartedAt)
	}
	choiceRaw, err := PrintList(taskIds, title)
	if err != nil {
		return nil, err
	}

	choice := strings.Split(choiceRaw, " | ")[0]
	var selected *congo.EcsTaskData
	for _, elem := range ecsTasks {
		if elem.Id == choice {
			selected = &elem
			break
		}
	}

	return selected, nil
}

func SelectElementFromEC2(instances *[]congo.Ec2Data, title string) (string, error) {
	var instancesDataMerge []string
	for _, i := range *instances {
		instancesDataMerge = append(
			instancesDataMerge,
			fmt.Sprintf("%s | %s | %s | %s | %s | %s",
				i.Name, i.Id, i.Platform, i.Type, i.Vpc, i.PrivateIp),
		)
	}

	choice, err := PrintList(instancesDataMerge, title)
	if err != nil {
		return "", err
	}

	if Exited {
		os.Exit(0)
	}

	selectedId := strings.Split(choice, " | ")[1]

	return selectedId, nil
}
