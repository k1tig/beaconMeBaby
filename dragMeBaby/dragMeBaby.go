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
type responseMsg struct{}

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
	stg       chan struct{}
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	keys      keyMap
	quitting  bool
	stage     stageTimes
}
type stageTimes struct {
	current   int
	beforestg int
	prestg    int
	stg       int
	yellow    float32
	green     bool
}

func (m model) listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		m.stage.beforestg = rand.Intn(stg1-4) + 4
		m.stage.prestg = rand.Intn(stg2-1) + 1
		m.stage.stg = rand.Intn(stg3-1) + 1
		m.stage.green = true
		for {
			switch {
			case m.stage.current == 0:
				go func() {
					time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
					m.stage.current++
					sub <- struct{}{}
				}()

			case m.stage.current == 1:
				go func() {
					time.Sleep(time.Second * time.Duration(m.stage.beforestg)) // nolint:gosec
					m.stage.current++
					sub <- struct{}{}
				}()
			case m.stage.current == 2:
				go func() {
					time.Sleep(time.Second * time.Duration(m.stage.prestg)) // nolint:gosec
					m.stage.current++
					sub <- struct{}{}
				}()
			case m.stage.current == 3:
				go func() {
					time.Sleep(time.Second * time.Duration(m.stage.stg)) // nolint:gosec
					m.stage.current++
					sub <- struct{}{}
				}()
			case m.stage.current == 4:
				go func() {
					time.Sleep(time.Second * time.Duration(m.stage.stg)) // nolint:gosec
					m.stage.current++
					sub <- struct{}{}
				}()

			}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),     // wait for activity
	)
}

// For staging a sequence of commands need to be sent for changing the lights. | before-stage (fault) pre-stage (fault) | staging
// (fault) | Yellow lights (fault) | Green Lights (start reaction timer) |

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Action): // location for action input / multi button
			switch {
			case m.stage.current == 0:
				m.stage.current++
			}
			return m, nil
		}
		// begin - wait, prestage - wait, stage- wait, yellow- wait, green
		return m, nil

	case responseMsg:

		switch {
		case m.stage.current == 1:
			m.stage.current++
		case m.stage.current == 2:
			m.stage.current++
		case m.stage.current == 3:
			m.stage.current++
		case m.stage.current == 4:
			m.stage.current = 0
		} // record external activity
		return m, waitForActivity(m.sub) // wait for next event
	default:
		return m, nil
	}
}
func (m model) View() string {
	s := fmt.Sprintf("\n Events received: %d\n\n Press any key to exit\n\n Level of stage: %d", m.responses, m.stage.current)
	if m.quitting {
		s += "\n"
	}
	return s
}
func main() {
	p := tea.NewProgram(model{
		stg:  make(chan struct{}),
		sub:  make(chan struct{}),
		keys: keys,
		stage: stageTimes{yellow: .400,
			current: 0,
		},
	})
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
