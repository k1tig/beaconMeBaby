package main

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type raceStatusMsg int

type keyMap struct {
	Action key.Binding
	Quit   key.Binding
}

var stg1 int = 6
var stg2 int = 4
var stg3 int = 4

var keys = keyMap{
	Action: key.NewBinding(
		key.WithKeys("g"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "Q"),
	),
}

type model struct {
	active   bool
	stg      chan int
	keys     keyMap
	quitting bool
	stage    stageTimes
}
type stageTimes struct {
	current   int
	beforestg int
	prestg    int
	stg       int
	yellow    float32
	green     bool
}

func (m model) setTimes() {
	s := m.stage
	s.beforestg = rand.Intn(stg1-4) + 4
	s.prestg = rand.Intn(stg2-1) + 1
	s.stg = rand.Intn(stg3-1) + 1
	s.green = true

}
func (m model) Init() tea.Cmd {
	return nil
}

func (m model) raceStatus(stg chan int) tea.Cmd {
	return func() tea.Msg {
		for {
			s := <-stg
			switch {
			case s == 1: //before pre-stage
				time.Sleep(time.Second * time.Duration(m.stage.beforestg)) // nolint:gosec
				m.stage.current++
				return raceStatusMsg(<-m.stg)
			case s == 2:
				time.Sleep(time.Second * time.Duration(m.stage.prestg)) // nolint:gosec
				m.stage.current++
				return raceStatusMsg(<-m.stg)
			case s == 3:
				time.Sleep(time.Second * time.Duration(m.stage.stg)) // nolint:gosec
				m.stage.current++
				return raceStatusMsg(<-m.stg)
			case s == 4:
				time.Sleep(time.Second * time.Duration(m.stage.stg)) // nolint:gosec
				m.stage.current++
			}
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//.7-1.3second stage to yellow
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Action): // location for action input / multi button
			switch {
			case !m.active:
				m.setTimes()
				m.active = true

				// ^^^^^ It doesn't update stage current until done go routine

				go m.raceStatus(m.stg)
				m.stg <- m.stage.current
				return m, nil
			}
			// begin - wait, prestage - wait, stage- wait, yellow- wait, green
			return m, nil

		default:
			return m, nil
		}
	case raceStatusMsg:
		s := m.stage.current
		switch {
		case s == 1:
			//Pre-stage lights
			m.stg <- s
		case s == 2:
			//Stage lights
			m.stg <- s
		case s == 3:
			m.stg <- s
			// Yellow lights

		}
	}
	return m, nil
}
func (m model) View() string {
	s := fmt.Sprintf("\n Current Stage: %d \n\n Press any key to exit\n\n ", m.stage.current)
	if m.quitting {
		s += "\n"
	}
	return s
}
func main() {
	p := tea.NewProgram(model{
		active: false,
		stg:    make(chan int),
		keys:   keys,
		stage:  stageTimes{yellow: .400, current: 1},
	})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
