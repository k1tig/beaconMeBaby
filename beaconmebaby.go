package main

// A simple program that counts down from 5 and then exits.
// https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	sub        chan struct{}
	keys       keyMap
	help       help.Model
	inputStyle lipgloss.Style
	itemStyle  lipgloss.Style
	lastKey    string
	quitting   bool
	position   int
	shift      int
	station    stations
	//	station    stations
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
		{k.Twenty, k.Seventeen, k.Fifteen, k.Twelve, k.Ten}, // first column
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
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		itemStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
	}
}
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	beacons := []stations{
		{"Madeira", "CS3B", -4, true},
		{"Argentina", "LU4AA", -3, true},
		{"Peru", "OA4B", -2, true},
		{"Venezuela", "YV5B", -1, true},
		{"United Nations", "4U1UN", 0, true},
		{"Canada", "VE8AT", 1, true},
		{"United States", "W6WX", 2, true},
		{"Hawaii", "KH6RS", 3, true},
		{"New Zealand", "ZL6B", 4, true},
		{"Australia", "VK6RBP", 5, true},
		{"Japan", "JA2IGY", 6, true},
		{"Russia", "RR9O", 7, true},
		{"Hong Kong", "VR2B", 8, true},
		{"Sri Lanka", "4S7B", 9, true},
		{"South Africa", "ZS6DN", 10, true},
		{"Kenya", "5Z4B", 11, true},
		{"Israel", "4X6TU", 12, true},
		{"Finland", "OH2B", 13, true},
		{"Madeira", "CS3B", 14, true},
		{"Argentina", "LU4AA", 15, true},
		{"Peru", "OA4B", 16, true},
		{"Venezuela", "YV5B", 17, true},
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
			m.shift = 0
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
		case key.Matches(msg, m.keys.Seventeen):
			m.lastKey = "17 meters"
			m.shift = -1
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
		case key.Matches(msg, m.keys.Fifteen):
			m.lastKey = "15 meters"
			m.shift = -2
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
		case key.Matches(msg, m.keys.Twelve):
			m.lastKey = "12 meters"
			m.shift = -3
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
		case key.Matches(msg, m.keys.Ten):
			m.lastKey = "10 meters"
			m.shift = -4
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
		return m, waitForActivity(m.sub)
	case responseMsg:
		a := startPosition()

		if a != nil {
			m.position = a.(int)
			for _, station := range beacons {
				if m.position+m.shift == station.id {
					m.station = station
				}
			}
			return m, waitForActivity(m.sub)
		}
	}
	return m, waitForActivity(m.sub)
}

type responseMsg struct{}

func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			//time.Sleep(time.Millisecond * 100) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func (m model) View() string {
	var status string
	var pos string

	if m.quitting {
		return "Bye!\n"
	}
	if m.lastKey == "" {
		status = "Selects a Band: "
	} else {
		i := m.position
		pos = strconv.Itoa(i + 1)
		callsign := m.inputStyle.Render("Station:  ") + m.itemStyle.Render(pos+") "+m.station.callsign)
		location := m.inputStyle.Render("Location: ") + m.itemStyle.Render(m.station.location)
		status = m.inputStyle.Render("Current station selected: ") + m.inputStyle.Render(m.lastKey) + "\n\n" + callsign + "\n" + location
	}
	helpView := m.help.View(m.keys)
	height := 8 - strings.Count(status, "\n") - strings.Count(helpView, "\n")
	return "\n" + status + "\n" + strings.Repeat("\n", height) + helpView
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
	func getPosition() tea.Msg {
		now := time.Now()
		if now.Second()%10 == 0 {
			totalSec := (now.Minute() * 60) + now.Second()
			if totalSec <= 180 {
				tSlot := totalSec / 10
				return tSlot
			} else {
				clean_time := totalSec % 180
				tSlot := clean_time / 10
				return tSlot
			}
		}
		return nil
	}
*/

func startPosition() tea.Msg {
	now := time.Now()
	totalSec := (now.Minute() * 60) + now.Second()
	if totalSec <= 180 {
		tSlot := totalSec / 10
		return tSlot
	} else {
		clean_time := now.Second() % 180
		tSlot := clean_time / 10
		return tSlot
	}
}
