package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func newModel() tea.Model {
	files, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	items := []list.Item{item("..")}
	for _, f := range files {
		items = append(items, item(f.Name()))
	}

	const defaultWidth = 0
	var listHeight = 14

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "What do you want for dinner?"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{list: l}
}
