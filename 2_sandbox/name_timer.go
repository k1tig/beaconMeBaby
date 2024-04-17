package main

// A simple program that counts down from 5 and then exits.

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//var num int
//var letter string

type model struct {
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
}

//type name struct {
//	position  int
//	character string
//}

type keyMap struct {
	Twenty    key.Binding
	Seventeen key.Binding
	Fifteen   key.Binding
	Twelve    key.Binding
	Ten       key.Binding
	Quit      key.Binding
	Help      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Twelve, k.Seventeen, k.Fifteen, k.Twelve, k.Ten}, // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Twenty: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("(1)", "20 meters"),
	),
	Seventeen: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("(2)", "17 meters"),
	),
	Fifteen: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("(3)", "15 meters"),
	),
	Twelve: key.NewBinding(
		key.WithKeys("4"),
		key.WithHelp("(4)", "12 meters"),
	),
	Ten: key.NewBinding(
		key.WithKeys("5"),
		key.WithHelp("(5)", "10 meters"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

func newModel() model {
	return model{
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("ff75b7")),
	}
}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Twenty):
			m.lastKey = "20 meters"
		case key.Matches(msg, m.keys.Seventeen):
			m.lastKey = "17 meters"
		case key.Matches(msg, m.keys.Fifteen):
			m.lastKey = "15 meters"
		case key.Matches(msg, m.keys.Twelve):
			m.lastKey = "12 meters"
		case key.Matches(msg, m.keys.Ten):
			m.lastKey = "10 meters"
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var status string
	if m.lastKey == "" {
		status = "Selects a Band"
	} else {
		status = "Current band selected: " + m.inputStyle.Render(m.lastKey)
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")

	return "\n" + status + strings.Repeat("\n", height) + helpView
}

// Main Function
func main() {
	if os.Getenv("HELP_DEBUG") != "" {
		f, err := tea.LogToFile("debug.log", "help")
		if err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		}
		defer f.Close() // nolint:errcheck
	}

	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}

/*
func getLetter() (num int, letter string) {
	now := time.Now()
	current_time := now.Second()
	num = 0
	letter = ""
	names := []name{
		{1, "A"},
		{2, "B"},
		{3, "C"},
		{4, "D"},
		{5, "E"},
		{6, "F"},
		{7, "G"},
		{8, "H"},
		{9, "I"},
		{0, "J"},
	}

	for _, pair := range names {
		if now.Second() != current_time {
			if now.Second() == pair.position {
				num = pair.position
				letter = pair.character
			} else if int(now.Second()/10) == pair.position {
				num = pair.position
				letter = pair.character
			}
		}
	}
	return num, letter
}
*/
