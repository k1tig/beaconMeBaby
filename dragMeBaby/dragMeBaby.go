package main

// A simple example that shows how to send activity to Bubble Tea in real-time
// through a channel.

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// A message used to indicate that activity has occurred. In the real world (for
// example, chat) this would contain actual data.
type responseMsg struct{}

type keyMap struct {
	Action key.Binding
	Quit   key.Binding
}

var keys = keyMap{
	Action: key.NewBinding(
		key.WithKeys("g"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "Q"),
	),
}

// Simulate a process that sends events at an irregular interval in real time.
// In this case, we'll send events on the channel at a random interval between
// 100 to 1000 milliseconds. As a command, Bubble Tea will run this
// asynchronously.
func listenForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		for {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(900)+100)) // nolint:gosec
			sub <- struct{}{}
		}
	}
}

// A command that waits for the activity on a channel.
func waitForActivity(sub chan struct{}) tea.Cmd {
	return func() tea.Msg {
		return responseMsg(<-sub)
	}
}

type model struct {
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
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

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity
	)
}

var stg1 int = 6
var stg2 int = 4
var stg3 int = 4

// For staging a sequence of commands need to be sent for changing the lights. | before-stage (fault) pre-stage (fault) | staging
// (fault) | Yellow lights (fault) | Green Lights (start reaction timer) |

func (m model) staging() {
	m.stage.beforestg = rand.Intn(stg1-4) + 4
	m.stage.prestg = rand.Intn(stg2-1) + 1
	m.stage.stg = rand.Intn(stg3-1) + 1
	m.stage.green = true
	switch {
	case m.stage.current == 0:
		time.Sleep(time.Second * time.Duration(m.stage.beforestg))
		return
	case m.stage.current == 1:
		time.Sleep(time.Second * time.Duration(m.stage.prestg))
		return
	case m.stage.current == 2:
		time.Sleep(time.Second * time.Duration(m.stage.stg))
		return
	case m.stage.current == 3:
		time.Sleep(time.Millisecond * time.Duration(m.stage.yellow))
	}

}

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
				m.staging()
				m.stage.current++
			case m.stage.current == 1: // prestage
				m.stage.current++
			case m.stage.current == 2: // stage
				m.stage.current++
			case m.stage.current == 3:
				// start timer
				m.stage.current = 4
			case m.stage.current == 4:
				//stop timer and print time

			}

			return m, nil
		}
		// begin - wait, prestage - wait, stage- wait, yellow- wait, green
	case responseMsg:
		m.responses++                    // record external activity
		return m, waitForActivity(m.sub) // wait for next event
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("\n %s Events received: %d\n\n Press any key to exit\n\n Level of stage: %d", m.spinner.View(), m.responses, m.stage.current)
	if m.quitting {
		s += "\n"
	}
	return s
}

func main() {
	p := tea.NewProgram(model{
		sub:     make(chan struct{}),
		spinner: spinner.New(),
		keys:    keys,
		stage: stageTimes{yellow: .400,
			current: 0,
		},
	})

	if _, err := p.Run(); err != nil {

		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
