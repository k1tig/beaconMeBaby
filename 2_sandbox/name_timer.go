package main

// A simple program that counts down from 5 and then exits.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

//var num int
//var letter string

type responseMsg struct{}

func listenForActivity(m model, sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			m.getPosition()
			time.Sleep(time.Millisecond * 100) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

type model struct {
	sub        chan struct{}
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	lastKey    string
	quitting   bool
	position   int
	station    stations
}

type stations struct {
	location string
	callsign string
	id       int
	activte  bool
}

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
		sub:        make(chan struct{}),
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("ff75b7")),
	}
}
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m, m.sub), // generate activity
		waitForActivity(m.sub),      // wait for activity
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	beacons := []stations{
		{"United Nations", "4U1UN", 1, true},
		{"Canada", "VE8AT", 2, true},
		{"United States", "W6WX", 3, true},
		{"Hawaii", "KH6RS", 4, true},
		{"New Zealand", "ZL6B", 5, true},
		{"Australia", "VK6RBP", 6, true},
		{"Japan", "JA2IGY", 7, true},
		{"Russia", "RR9O", 8, true},
		{"Hong Kong", "VR2B", 9, true},
		{"Sri Lanka", "4S7B", 10, true},
		{"South Africa", "ZS6DN", 11, true},
		{"Kenya", "5Z4B", 12, true},
		{"Israel", "4X6TU", 13, true},
		{"Finland", "OH2B", 14, true},
		{"Madeira", "CS3B", 15, true},
		{"Argentina", "LU4AA", 16, true},
		{"Peru", "OA4B", 17, true},
		{"Venezuela", "YV5B", 18, true},
	}

	//This needs to be done right so that there is a constant running even that sends a message when changed.

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:

		switch {
		case key.Matches(msg, m.keys.Twenty):
			m.lastKey = "20 meters"
			for _, beacon := range beacons {
				if m.position == beacon.id {
					m.station = beacon
				}
			}

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
	case responseMsg:
		return m, waitForActivity(m.sub)

	}

	return m, m.getPosition
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var status string

	if m.lastKey == "" {
		status = "Selects a Band"
	} else {
		status = "Current station selected: " + m.inputStyle.Render(m.lastKey)
	}

	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")
	callsign := m.station.callsign
	return "\n" + status + "\n" + callsign + strings.Repeat("\n", height) + helpView
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

func (m model) getPosition() tea.Msg {
	now := time.Now()

	if now.Second()%10 != 0 {
		time.Sleep(time.Second * 1)
	} else {
		totalSec := (now.Minute() * 60) + now.Second()
		if totalSec <= 180 {
			tSlot := totalSec / 10
			m.position = tSlot
		} else {
			clean_time := totalSec % 180
			tSlot := clean_time / 10
			m.position = tSlot
		}
	}
	return m
}
