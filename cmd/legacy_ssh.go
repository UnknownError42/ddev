package cmd

import (
	"fmt"
	"log"
	"path"

	"github.com/drud/bootstrap/cli/local"
	"github.com/drud/drud-go/utils"
	"github.com/spf13/cobra"
)

// LegacySSHCmd represents the ssh command.
var LegacySSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "SSH to an app container.",
	Long:  `Connects user to the running container.`,
	Run: func(cmd *cobra.Command, args []string) {

		if activeApp == "" {
			log.Fatalln("Must set app flag to dentoe which app you want to work with.")
		}

		app := local.LegacyApp{
			Name:        activeApp,
			Environment: activeDeploy,
		}

		nameContainer := fmt.Sprintf("%s-%s", app.ContainerName(), serviceType)

		if !utils.IsRunning(nameContainer) {
			log.Fatal("App not running locally. Try `drud legacy add`.")
		}

		if !app.ComposeFileExists() {
			log.Fatalln("No docker-compose yaml for this site. Try `drud legacy add`.")
		}

		err := utils.DockerCompose(
			"-f", path.Join(app.AbsPath(), "docker-compose.yaml"),
			"exec",
			nameContainer,
			"bash",
		)
		if err != nil {
			log.Fatalln(err)
		}

	},
}

func init() {
	LegacySSHCmd.Flags().StringVarP(&serviceType, "service", "s", "web", "Which service to send the command to. [web, db]")
	LegacyCmd.AddCommand(LegacySSHCmd)
}
