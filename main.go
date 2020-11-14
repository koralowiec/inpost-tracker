package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/enescakir/emoji"
	"github.com/goodsign/monday"
	"gitlab.com/koralowiec/inpost-track/data"
)

var filePath string

type model struct {
	trackingNumbers []data.TrackingNumber
	trackingInfo    map[string]data.TrackingResponse
	showTrackInfo   map[string]bool
	cursor          int
	addNumberInput  textinput.Model
	addingNumber    bool
	err             error
}

func initalModel() model {
	inputModel := textinput.NewModel()
	inputModel.Placeholder = "Nowy numer przesyłki"
	inputModel.CharLimit = 156
	inputModel.Width = 24
	inputModel.Prompt = " " + emoji.Plus.String() + " "

	return model{
		addNumberInput: inputModel,
		trackingInfo:   make(map[string]data.TrackingResponse),
		showTrackInfo:  make(map[string]bool),
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

func getTrackingInformationFromAPI(trackingNumber data.TrackingNumber) tea.Cmd {
	return func() tea.Msg {
		trackRes, err := data.GetTrackingInfo(string(trackingNumber))
		if err != nil {
			return errMsg{err}
		}
		return trackingInfoMsg(*trackRes)
	}
}

type trackingInfoMsg data.TrackingResponse

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

	case trackingInfoMsg:
		number := msg.TrackingNumber
		m.trackingInfo[number] = data.TrackingResponse(msg)
		m.showTrackInfo[number] = true
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
			case "esc":
				m.addNumberInput.SetValue("")
				m.addingNumber = false
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
				} else {
					i := m.cursor - 1
					trackNum := m.trackingNumbers[i]
					trackNumString := string(trackNum)
					if _, ok := m.trackingInfo[trackNumString]; ok {
						if show, ok := m.showTrackInfo[trackNumString]; ok {
							m.showTrackInfo[trackNumString] = !show
						}
					} else {
						return m, getTrackingInformationFromAPI(m.trackingNumbers[i])
					}

					// } else {
					// 	i := m.cursor - 1
					// 	return m, getTrackingInformationFromAPI(m.trackingNumbers[i])
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
	var s string

	if m.err != nil {
		return fmt.Sprintf("\nWystąpił problem: %v \n\n", m.err)
	}

	if m.addingNumber {
		s += fmt.Sprintf(m.addNumberInput.View()) + "\n"
	} else {
		prefix := "  "
		if m.cursor == 0 {
			prefix = emoji.NewButton.String()
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

		trackInfo, ok := m.trackingInfo[string(trackNumber)]
		if ok {
			for _, status := range trackInfo.TrackingDetails {
				t := monday.Format(status.DateTime, "Mon 15:04 02.01.2006", monday.LocalePlPL)
				s += fmt.Sprintf(" %s %s %s %s \n", "  ", "  ", t, status.Status.Title)
			}
		}
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
