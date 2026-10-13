package main

import (
	"math/rand/v2"
	"strconv"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ModelState int

const (
	UsernameInput = iota
	Game
)

type model struct {
	*app
	state         ModelState
	username      string
	width         int
	height        int
	err           error
	usernameInput textinput.Model
	guessInput    textinput.Model
	answer        int
	guess         int
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func NewModel(username string) model {
	ui := textinput.New()
	ui.Placeholder = "Username..."
	ui.Cursor()
	ui.Focus()
	ui.Prompt = " "

	ui.SetValue(username)
	ui.CharLimit = 18

	ui.SetWidth(20)

	gi := textinput.New()
	gi.Placeholder = "Guess..."
	gi.Cursor()
	gi.Focus()
	gi.Prompt = " "

	gi.CharLimit = 7

	gi.SetWidth(8)

	return model{
		usernameInput: ui,
		guessInput:    gi,
		err:           nil,
		username:      username,
		state:         Game,
		answer:        rand.IntN(1000000),
		guess:         -1,
	}

}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		uiCmd tea.Cmd
		giCmd tea.Cmd
	)

	m.usernameInput, uiCmd = m.usernameInput.Update(msg)
	m.guessInput, giCmd = m.guessInput.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Key().Mod {
		case tea.ModCtrl: // We're only interested in ctrl+<key>
			switch msg.Key().Code {
			case 'c':
				return m, tea.Quit
			}
		}
		switch msg.Key().Code {
		case tea.KeyEnter:
			switch m.state {
			case UsernameInput:
				m.state = Game
			case Game:
				i, err := strconv.Atoi(m.guessInput.Value())
				if err == nil && i > -1 {
					m.guess = i
					m.guessInput.SetValue("")
				}
			}
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(uiCmd, giCmd)
}

func (m model) View() tea.View {
	var v tea.View
	switch m.state {
	case UsernameInput:
		v = tea.NewView(
			lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
				lipgloss.JoinVertical(
					lipgloss.Center,
					BANNER,
					"",
					lipgloss.JoinVertical(
						lipgloss.Center,
						GreenTextStyle.Render("Please Enter a Username:"),
						TextBoxStyle.Render(m.usernameInput.View())),
				),
			))
	case Game:
		var isLower bool = m.guess < m.answer
		var hint string
		if isLower {
			hint = "HIGHER!"
		} else {
			hint = "LOWER!"
		}
		if m.guess == m.answer {
			hint = "Correct!"
		}
		if m.guess == -1 {
			hint = ""
		}
		v = tea.NewView(lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center,
				BANNER,
				"",
				GreenTextStyle.Render("Guess the number! 0-1M"),
				"",
				GreenTextStyle.Render(hint),
				"",
				TextBoxStyle.Render(m.guessInput.View()),
			),
		))
	}
	v.AltScreen = true
	return v
}
