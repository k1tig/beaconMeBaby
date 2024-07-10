package main

// A simple program that counts down from 5 and then exits.
// https://github.com/charmbracelet/bubbletea/blob/master/examples/realtime/main.go
import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	stopwatch time.Duration

	sub      chan struct{}
	keys     keyMap
	help     help.Model
	active   bool
	quitting bool
	stg      int
	stgT     times
	timer    time.Time //	station    stations
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
	Action key.Binding
}

var keys = keyMap{
	Action: key.NewBinding(
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
		stgT: times{preStg: 2, fullStg: 2, Yellow: 1.2, Green: .400},
	}
}
func (m model) Init() tea.Cmd {
	return tea.Batch(
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity

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
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func (m model) View() string {
	s := "Elapsed Time: " + m.stopwatch.String() + "\nCurrent Stage: " + strconv.Itoa(m.stg)
	return s
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
		case key.Matches(msg, m.keys.Action):
			if m.stg == 4 {
				x := time.Now()
				m.stopwatch = x.Sub(m.timer)
				m.stg++
				m.timer = time.Now()
				m.active = true

			}
		}
		return m, waitForActivity(m.sub)
	case responseMsg:
		switch {
		case !m.active && m.stg == 0:
			now := time.Now()
			m.timer = now
			m.active = true
			m.keys.Action.SetEnabled(true)
		case m.active:
			switch {
			case m.stg == 0:
				x := m.timer.Add(time.Millisecond * time.Duration(1000*m.stgT.preStg))
				current := time.Now()
				if current.After(x) {
					m.stg++
					m.timer = current

				}

			case m.stg == 1:
				x := m.timer.Add(time.Millisecond * time.Duration(1000*m.stgT.fullStg))
				current := time.Now()
				if current.After(x) {
					m.stg++
					m.timer = current
				}

			case m.stg == 2:
				x := m.timer.Add(time.Millisecond * time.Duration(1000*m.stgT.Yellow))
				current := time.Now()
				if current.After(x) {
					m.stg++
					m.timer = current
				}
			case m.stg == 3:
				x := m.timer.Add(time.Millisecond * time.Duration(1000*m.stgT.Green))
				current := time.Now()
				if current.After(x) {
					m.stg++
					m.timer = current

				}
			case m.stg == 4:
				m.timer = time.Now()
				m.active = false

			case m.stg == 5:
				x := m.timer.Add(time.Second * 5)
				current := time.Now()
				if current.After(x) {
					m.stg = 0
					m.active = false

				}

			}

		}
	}
	return m, waitForActivity(m.sub)
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
