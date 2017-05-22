package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init will initialize the synchronization.",
	Long: `Initialize (dotsync init) will clone your dotfiles repository into your user directory,
under  "~/.dotfiles" by default and:

  * Files in "/copy" are copied unto "~/";
  * Files in "/link" are symlinked into "~/".
  * You will be prompted to run "init" scripts and setups.append
  	The installer will select only the relevant init processes based on your OS.

The "~/.backups" directory gets created when necessary. Any files in ~/ that would have been overwritten by files in /copy or /link get backed up there.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("init called")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)

	initCmd.Flags().StringP("dotfiles-dir", "d", "~/.dotfiles", "Directory used to clone the dotfiles")
	initCmd.Flags().StringP("backup-dir", "b", "~/.backups", "Directory used to storage the backups")
}
