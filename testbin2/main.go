package main

// A simple program that counts down from 5 and then exits.
// https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
import (
	"fmt"
	"math/rand"
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
	quitting   bool
	pos        int
	stg        int
	stgTimes   stgTimelst

	//	station    stations
}

type stgTimelst struct {
	prestg int
	stg    int
	yellow int // ,7 - 1.3 sec
	green  float32
	active bool
}

type keyMap struct {
	Twenty key.Binding
	Quit   key.Binding
	Help   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Twenty},       // first column
		{k.Help, k.Quit}, // second column
	}
}

var keys = keyMap{
	Twenty: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("(g)", "Action"),
	),
}

func newModel() model {
	return model{
		sub:        make(chan struct{}),
		keys:       keys,
		help:       help.New(),
		inputStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("6")),
		itemStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		stgTimes:   stgTimelst{green: .400, active: false},
		stg:        0,
		pos:        1,
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
		case key.Matches(msg, m.keys.Twenty):
			if !m.stgTimes.active {
				m.setTimes()
			}
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
		return m, waitForActivity(m.sub)
	case responseMsg:
		switch m.stgTimes.active {
		case m.stg != m.pos:
			switch {
			case m.stg == 0:
				stgDelay(m.stgTimes.prestg)
				m.stg++
			case m.stg == 1:
				stgDelay(m.stgTimes.stg)
				m.stg++
			case m.stg == 2:
				stgDelay(m.stgTimes.yellow)
				m.stg++
			case m.stg == 3:
				stgDelay(int(m.stgTimes.green))
				m.stg++
			case m.stg == 4:
				m.stg = 0
				m.pos = 1
				m.stgTimes.active = false

			}
			return m, waitForActivity(m.sub)

		case m.stg == m.pos:
			m.pos++

		}
		return m, waitForActivity(m.sub)
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

func stgDelay(x int) {
	time.Sleep(time.Second * time.Duration(x))
}

// presets for stage times on drag tree
var stg1 = 6
var stg2 = 3
var stg3 = 130

func (m model) setTimes() tea.Model {
	s := m.stgTimes
	s.prestg = rand.Intn(stg1-4) + 4
	s.stg = rand.Intn(stg2-1) + 1
	s.yellow = (rand.Intn(stg3-70) + 70) / 10
	s.active = true
	return m
}

func (m model) View() string {
	var status string
	var stg string

	if m.quitting {
		return "Bye!\n"

	} else {
		i := m.stg
		stg = strconv.Itoa(i)

		status = m.inputStyle.Render("Stage:") + m.inputStyle.Render(stg)
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
