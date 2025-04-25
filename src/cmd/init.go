package cmd

import (
	"fmt"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/ErnieBernie10/simplecloud/internal"
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

type Config struct {
	Domain string `toml:"domain"`
	Email  string `toml:"email"`
}
type TomlConfig struct {
	Config Config `toml:"Config"`
}

func saveResults(m model) (tea.Model, tea.Cmd) {
	// Write toml file with results
	m.domain = m.inputs[domain].Value()
	m.email = m.inputs[email].Value()
	config := TomlConfig{
		Config: Config{
			Domain: m.domain,
			Email:  m.email,
		},
	}
	content, err := toml.Marshal(config)
	if err != nil {
		m.err = err
		return m, nil
	}
	if err := os.MkdirAll(internal.TargetDir, 0755); err != nil {
		m.err = err
		return m, nil
	}
	err = os.WriteFile(fmt.Sprintf("%s/config.toml", internal.TargetDir), content, 0644)
	if err != nil {
		m.err = err
		return m, nil
	}

	if err = createInitDockerCompose(); err != nil {
		m.err = err
		return m, nil
	}
	if err = createInitEnv(config); err != nil {
		m.err = err
		return m, nil
	}
	fmt.Println("Config saved successfully")

	return m, tea.Quit
}

func createInitEnv(config TomlConfig) error {
	env, err := internal.Opt.ReadFile("traefik/.env")
	if err != nil {
		return err
	}
	tmpl, err := template.New("env").Parse(string(env))
	if err != nil {
		return err
	}
	envFile, err := os.Create(fmt.Sprintf("%s/.env", internal.TargetDir))
	if err != nil {
		return err
	}
	defer envFile.Close()
	err = tmpl.Execute(envFile, map[string]string{
		"Domain": config.Config.Domain,
		"Email":  config.Config.Email,
	})
	if err != nil {
		return err
	}
	return nil
}

func createInitDockerCompose() error {
	traefik, err := internal.Opt.ReadFile("traefik/docker-compose.yml")
	if err != nil {
		return err
	}
	err = os.WriteFile(fmt.Sprintf("%s/docker-compose.yml", internal.TargetDir), traefik, 0644)
	if err != nil {
		return err
	}
	return nil
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
