package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/ErnieBernie10/simplecloud/src/internal"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

type model struct {
	domain  string `toml:"domain"`
	email   string `toml:"email"`
	focused int
	err     error
	inputs  []textinput.Model
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize tool",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[domain] = textinput.New()
	inputs[domain].Placeholder = "Domain"
	inputs[domain].Focus()
	inputs[domain].Prompt = "Domain: "
	inputs[domain].CharLimit = 156
	inputs[domain].Width = 20

	inputs[email] = textinput.New()
	inputs[email].Placeholder = "Email"
	inputs[email].Prompt = "Email: "
	inputs[email].CharLimit = 156
	inputs[email].Width = 20

	return model{
		focused: 0,
		inputs:  inputs,
		err:     nil,
		domain:  "",
		email:   "",
	}
}

// Init implements tea.Model.
func (m model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return saveResults(m)
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func saveResults(m model) (tea.Model, tea.Cmd) {
	// Write toml file with results
	m.domain = m.inputs[domain].Value()
	m.email = m.inputs[email].Value()
	config := internal.TomlConfig{
		Config: internal.Config{
			Domain: m.domain,
			Email:  m.email,
		},
	}

	run := internal.NewRunContext(internal.TargetDir, context.Background())
	err := run.Bootstrap(config)
	if err != nil {
		m.err = fmt.Errorf("failed to save config: %w", err)
		return m, tea.Quit
	}

	fmt.Println("Config saved successfully")

	return m, tea.Quit
}

const (
	domain = iota
	email
)

// View implements tea.Model.
func (m model) View() string {

	if m.err != nil {
		return m.err.Error()
	}
	return fmt.Sprintf(
		`
 %s
 %s
`,
		m.inputs[domain].View(),
		m.inputs[email].View(),
	) + "\n"

}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
