package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

var staging float32

const prestage = 10
const fullstage = 7
const raceyellow = 1.3
const racegreen = .4

const runprogram = true

func racetime(stage int) {
	p1 := time.Duration(rand.Intn(stage+1) + 1)
	time.Sleep(p1 * time.Second)
}

type model struct {
	timer    timer.Model
	keymap   keymap
	help     help.Model
	quitting bool
	rStatus  raceStatus
}

type raceStatus struct {
	begin    bool
	preStage bool
	stage    bool
	yellow   bool
	green    bool
	end      bool
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	quit  key.Binding
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd

	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, m.keymap.start):
			switch {
			case m.rStatus.begin != true:
				m.rStatus.begin = true
				racetime(prestage)
				return m, nil

			case m.rStatus.preStage != true:
				m.rStatus.preStage = true
				racetime(fullstage)
				return m, nil

			case m.rStatus.yellow != true:
				m.rStatus.yellow = true
				time.Sleep((raceyellow * 1000) * time.Millisecond)
				return m, nil

			case m.rStatus.green != true:
				m.rStatus.green = true
				time.Sleep((racegreen * 1000) * time.Millisecond)
				return m, nil // Send CMD to start timer. Above Code shouldnt be bound to start key.

			}
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,

		m.keymap.quit,
	})
}

func (m model) View() string {
	// For a more detailed timer view you could read m.timer.Timeout to get
	// the remaining time as a time.Duration and skip calling m.timer.View()
	// entirely.
	s := m.timer.View()

	if m.timer.Timedout() {
		s = "All done!"
	}
	s += "\n"
	if !m.quitting {
		s = "Exiting in " + s
		s += m.helpView()
	}
	return s
}

func main() {
	r1 := raceStatus{
		false, false, false, false,
	}
	fmt.Printf("Staring Race\n")
	for runprogram == true {

		r1.current = true
		// Press start key

		r1.foulOn = true
		for r1.foulOn {

			p1 := time.Duration(rand.Intn(prestage+1) + 1)
			time.Sleep(p1 * time.Second)
			fmt.Println("Prestage... ")

			s1 := time.Duration(rand.Intn(fullstage+1) + 1)
			time.Sleep(s1 * time.Second)
			fmt.Printf("Fullstage...\n")

			time.Sleep(time.Duration(raceyellow*1000) * time.Millisecond)
			fmt.Printf("oo|oo\noo|oo\n 0|0\n 0|0\n 0|0\n")

			time.Sleep(time.Duration(racegreen*1000) * time.Millisecond)
			fmt.Println(" X|X")

			r1.foulOn = false
			r1.greenlight = true
		}

		for r1.greenlight {
			race1 := time.Duration(rand.Intn(100+1) + 20)
			time.Sleep((race1 * time.Millisecond))
			race2 := float32(race1)
			fmt.Printf("Race Results: 0.%g\n ", race2)
			r1.greenlight = false

		}

		r1.endrace = true
		fmt.Printf("End Race\n\n")
		time.Sleep(5 * time.Second)
	}

}
