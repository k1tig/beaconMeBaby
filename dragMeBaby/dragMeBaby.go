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
	levels    int
	sub       chan struct{} // where we'll receive activity notifications
	responses int           // how many responses we've received
	spinner   spinner.Model
	keys      keyMap
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		listenForActivity(m.sub), // generate activity
		waitForActivity(m.sub),   // wait for activity
	)
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
			case m.levels == 0: // make channel to listen for flase start
				//time.Sleep(time.Second * 4)
				m.levels++
			case m.levels == 1: // prestage
				m.levels++
			case m.levels == 2: // stage
				m.levels++
			case m.levels > -3:
				m.levels += 2
			}

			return m, nil
		}

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
	s := fmt.Sprintf("\n %s Events received: %d\n\n Press any key to exit\n\n Level of stage: %d", m.spinner.View(), m.responses, m.levels)
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
		levels:  1,
	})

	if _, err := p.Run(); err != nil {

		fmt.Println("could not start program:", err)
		os.Exit(1)
	}
}
