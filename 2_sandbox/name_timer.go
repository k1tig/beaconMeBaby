package main

// A simple program that counts down from 5 and then exits.

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Log to a file. Useful in debugging since you can't really log to stdout.
	// Not required.
	logfilePath := os.Getenv("BUBBLETEA_LOG")
	if logfilePath != "" {
		if _, err := tea.LogToFile(logfilePath, "simple"); err != nil {
			log.Fatal(err)
		}
	}

	// Initialize our program
	p := tea.NewProgram(model(6))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

var num int
var letter string

type model int
type name struct {
	position  int
	character string
}

func sleep() {
	time.Sleep(time.Second * 1)
}

func getLetter() (num int, letter string) {
	now := time.Now()
	current_time := now.Second()
	num = 0
	letter = ""
	names := []name{
		{1, "A"},
		{2, "B"},
		{3, "C"},
		{4, "D"},
		{5, "E"},
		{6, "F"},
		{7, "G"},
		{8, "H"},
		{9, "I"},
		{0, "J"},
	}

	for _, pair := range names {
		if now.Second() != current_time {
			if now.Second() == pair.position {
				num = pair.position
				letter = pair.character
			} else if int(now.Second()/10) == pair.position {
				num = pair.position
				letter = pair.character
			}
		}
	}
	return num, letter
}

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit
	case tickMsg:
		m--
		if m <= 0 {
			return m, tea.Quit
		}
		return m, tick
	}
	return m, nil
}

func (m model) View() string {
	return fmt.Sprintf("counter %v and Name %v\n", num, letter)
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
