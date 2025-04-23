package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/ErnieBernie10/simplecloud/internal"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [app]",
	Short: "Install a Docker Compose app",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app := args[0]
		err := installApp(app)
		if err != nil {
			fmt.Printf("Failed to install app %s: %v\n", app, err)
			os.Exit(1)
		} else {
			fmt.Printf("App %s installed successfully!\n", app)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

func getConfig() (TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(fmt.Sprintf("%s/config.toml", internal.TargetDir), &config); err != nil {
		return config, err
	}
	return config, nil
}

func installApp(appName string) error {
	app, err := internal.GetApp(appName)
	if err != nil {
		return err
	}
	composeFile := fmt.Sprintf("%s/docker-compose.yml", app.TargetDir)
	if err := os.MkdirAll(app.TargetDir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(composeFile, []byte(app.DockerCompose), 0644); err != nil {
		return err
	}
	tmpl, err := template.New("env").Parse(app.EnvTemplate)
	if err != nil {
		return err
	}
	config, err := getConfig()
	if err != nil {
		return err
	}

	envFile, err := os.Create(fmt.Sprintf("%s/.env", app.TargetDir))
	if err != nil {
		return err
	}
	defer envFile.Close()
	// TODO: Accept subdomain as parameter
	// TODO: Accept port as parameter
	err = tmpl.Execute(envFile, map[string]string{
		"Domain": config.Config.Domain,
		"Email":  config.Config.Email,
	})
	if err != nil {
		return err
	}

	if !internal.IsMasterRunning() {
		return fmt.Errorf("master is not running. Please start the service before installing an app")
	}

	cmd := exec.Command("docker", "compose", "-f", composeFile, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
