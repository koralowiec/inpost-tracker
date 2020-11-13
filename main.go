package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/koralowiec/inpost-track/data"
)

const filePath = "/home/arek/.inpost-track/saved.json"

type model struct {
	trackingNumbers []data.TrackingNumber
	cursor          int
	selected        map[int]struct{}
	addNumberInput  textinput.Model
	addingNumber    bool
	err             error
}

func initalModel() model {
	inputModel := textinput.NewModel()
	inputModel.Placeholder = "Nowy numer przesyłki"
	inputModel.CharLimit = 156
	inputModel.Width = 20

	return model{
		addNumberInput: inputModel,
	}
}

func getSavedTrackingNumbers() tea.Msg {
	fc := data.LoadFileContent(filePath)
	return trackingNumbersMsg(fc.TrackingNumbers)
}

type trackingNumbersMsg []data.TrackingNumber

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return getSavedTrackingNumbers
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case trackingNumbersMsg:
		m.trackingNumbers = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.trackingNumbers) {
				m.cursor++
			}
		case "enter":
			if m.cursor == 0 {
				m.addingNumber = !m.addingNumber
				if m.addingNumber {
					m.addNumberInput.Focus()
					textinput.Blink()
				}
			}
		}
	}

	var cmd tea.Cmd
	m.addNumberInput, cmd = m.addNumberInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWystąpił problem: %v\n\n", m.err)
	}

	var s string

	if m.addingNumber {
		s += fmt.Sprintf(m.addNumberInput.View()) + "\n"
	} else {
		prefix := " "
		if m.cursor == 0 {
			prefix = ">"
		}
		s += fmt.Sprintf(" %s %s\n", prefix, "Dodaj nowy numer przesyłki")
	}

	for i, trackNumber := range m.trackingNumbers {
		i += 1
		prefix := " "
		if m.cursor == i {
			prefix = ">"
		}
		s += fmt.Sprintf(" %s %s\n", prefix, trackNumber)
	}

	return "\n" + s + "\n\n"
}

func main() {
	program := tea.NewProgram(initalModel())
	if err := program.Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
