package main

// A simple program that counts down from 5 and then exits.
// https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	sub      chan struct{}
	keys     keyMap
	help     help.Model
	quitting bool
	stg      int
	stgT     times
	color    string

	//	station    stations
}

type times struct {
	preStg  int
	fullStg int
	Yellow  float32
	Green   float32
}

type keyMap struct {
	Twenty key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Twenty: key.NewBinding(
		key.WithKeys("g"),
		key.WithHelp("(g)", "Action"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("(q)", "Quit"),
	),
}

func newModel() model {
	return model{
		sub:  make(chan struct{}),
		keys: keys,
		help: help.New(),
		stg:  0,
		stgT: times{preStg: 4, fullStg: 3, Yellow: 1.2, Green: .400},
	}
}
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub), // generate activity
		m.waitForActivity(m.sub), // wait for activity

	)
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
func (m model) waitForActivity(sub chan struct{}) tea.Cmd {
	switch {
	case m.stg == 0:
		time.Sleep(time.Second * time.Duration(m.stgT.preStg))

	case m.stg == 1:
		time.Sleep(time.Second * time.Duration(m.stgT.fullStg))

	case m.stg == 2:
		time.Sleep(time.Millisecond * time.Duration(m.stgT.Yellow*1000))

	case m.stg == 3:
		time.Sleep(time.Millisecond * time.Duration(m.stgT.Green*1000))

	case m.stg == 4:
		time.Sleep(time.Second * time.Duration(5))

	}
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}
func (m model) View() string {
	return fmt.Sprintf("Current stage:%d\nStatus:%v", m.stg, m.color)
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
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		}
		return m, m.waitForActivity(m.sub)
	case responseMsg:
		switch {
		case m.stg == 0:
			m.stg++
			m.color = "Prestage"
			return m, m.waitForActivity(m.sub)

		case m.stg == 1:
			m.stg++
			m.color = "Stage"
			return m, m.waitForActivity(m.sub)

		case m.stg == 2:
			m.stg++
			m.color = "Yellow"
			return m, m.waitForActivity(m.sub)

		case m.stg == 3:
			m.stg++
			m.color = "Green"
			return m, m.waitForActivity(m.sub)

		case m.stg == 4:
			time.Sleep(time.Second * time.Duration(5))
			return m, m.waitForActivity(m.sub)

		}
	}
	return m, m.waitForActivity(m.sub)
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

//  ____________
// |(oo)=||=(oo)|
// |(oo)=||=(oo)|
//   ==========
//  |(O)=||=(O)|
//  |(O)=||=(O)|
//  |(O)=||=(O)|
//  |====||====|
//  |(O)=||=(O)|
//   ==========
//      ||||
//      ||||
//      ||||
//      ||||
//      ||||
//      ||||
// --------------
