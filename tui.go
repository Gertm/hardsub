package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type tuiModel struct {
	cursor     int
	startFrame string
	endFrame   string
}

func initialModel() tuiModel {
	return tuiModel{
		cursor:     0,
		startFrame: "",
		endFrame:   "",
	}
}

func (m tuiModel) Init() tea.Cmd {
	return nil
}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msgt := msg.(type) {
	case tea.KeyMsg:
		switch msgt.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor -= 1
			}
		case "down", "j":
			if m.cursor < 1 {
				m.cursor += 1
			}
		}
	}
	return m, nil
}

func (m tuiModel) View() string {
	sb := strings.Builder{}
	sb.WriteString("Hardsub frame cutting dialog\n\n")
	if m.cursor == 0 {
		sb.WriteString("> ")
	} else {
		sb.WriteString("  ")
	}
	sb.WriteString("Start frame: \n\n")
	if m.cursor == 1 {
		sb.WriteString("> ")
	} else {
		sb.WriteString("  ")
	}
	sb.WriteString("End frame: \n\n")
	sb.WriteString("q or ctrl-c to quit, arrows or j/k to move up and down\n")
	return sb.String()
}

func startTui() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
