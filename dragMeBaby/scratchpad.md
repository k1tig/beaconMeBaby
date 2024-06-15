package main

// A simple program that counts down from 5 and then exits.
// https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"math/rand/v2"

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
	countdown bool
	raceStatus bool
}

type keyMap struct {
	Start key.Binding
	Go    key.Binding
	Quit  key.Binding
	Help  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Start, k.Go},  // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Start: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("(s)", "Start"),
	),

	Go: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("(g)", "Go"),
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

	//This needs to be done right so that there is a constant running even that sends a message when changed.
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// If we set a width on the help menu it can gracefully truncate
		// its view as needed.
		m.help.Width = msg.Width

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Start):
			if m.countdown == false{
				pS = rand
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
		if m.position != a.(int) {
			if a != nil {
				m.position = a.(int)
				for _, station := range beacons {
					if m.position+m.shift == station.id {
						m.station = station

					}
				}
				//time.Sleep(time.Millisecond * 1000)
				return m, waitForActivity(m.sub)
			}
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

func startPosition() tea.Msg {
	now := time.Now()
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
