package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enescakir/emoji"
	"gitlab.com/koralowiec/inpost-track/data"
)

var filePath string

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
	fc, err := data.LoadFileContent(filePath)
	if err != nil {
		return errMsg{err}
	}
	return trackingNumbersMsg(fc.TrackingNumbers)
}

type trackingNumbersMsg []data.TrackingNumber

func appendNewTrackingNumber(filePath string, trackNum string) tea.Cmd {
	trackNumber := data.TrackingNumber(trackNum)
	return func() tea.Msg {
		trackingNumbers, err := data.AppendTrackingNumber(filePath, trackNumber)
		if err != nil {
			return errMsg{err}
		}
		return trackingNumbersMsg(trackingNumbers)
	}
}

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

	case errMsg:
		m.err = msg
		return m, nil

	case tea.KeyMsg:
		// Update that should happen always when that key/keys are pressed
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.cursor == 0 {
				m.addingNumber = !m.addingNumber
			}
		}

		// Update depends on whether user is adding a new package number
		if m.addingNumber {
			switch msg.String() {
			case "enter":
				if m.cursor == 0 {
					m.addNumberInput.Focus()
					textinput.Blink()
				}

			}
		} else {
			switch msg.String() {
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
					newNumber := m.addNumberInput.Value()
					m.addNumberInput.SetValue("")
					return m, appendNewTrackingNumber(filePath, newNumber)
				}
			}

		}
	}

	var cmd tea.Cmd
	if m.addingNumber {
		m.addNumberInput, cmd = m.addNumberInput.Update(msg)
	}
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
		prefix := "  "
		if m.cursor == 0 {
			prefix = emoji.Package.String()
		}
		s += fmt.Sprintf(" %s %s\n", prefix, "Dodaj nowy numer przesyłki")
	}

	for i, trackNumber := range m.trackingNumbers {
		i += 1
		prefix := "  "
		if m.cursor == i {
			prefix = emoji.Package.String()
		}
		s += fmt.Sprintf(" %s %s\n", prefix, trackNumber)
	}

	return "\n" + s + "\n\n"
}

func main() {
	var err error
	filePath, err = data.GetContentFilePath()
	if err != nil {
		fmt.Printf("Nie udało się uzyskać ścieżki do pliku, błąd: %+v\n", err)
	}

	program := tea.NewProgram(initalModel())
	if err := program.Start(); err != nil {
		fmt.Printf("O nie, góra lodowa: %v tszszszszs\n", err)
		os.Exit(1)
	}
}
