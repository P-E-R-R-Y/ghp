package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate or refresh GitHub CLI token with optional scopes",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		permission, _ := cmd.Flags().GetString("permission")

		if host == "" {
			host = "github.com"
		}

		var ghCmd *exec.Cmd
		if permission != "" {
			fmt.Printf("Refreshing auth for host %s with permission %s...\n", host, permission)
			ghCmd = exec.Command("gh", "auth", "refresh", "-h", host, "-s", permission)
		} else {
			fmt.Printf("Checking auth for host %s...\n", host)
			ghCmd = exec.Command("gh", "auth", "status", "-h", host)
		}

		ghCmd.Stdout = os.Stdout
		ghCmd.Stderr = os.Stderr

		err := ghCmd.Run()
		if err != nil {
			return fmt.Errorf("gh auth command failed: %w", err)
		}

		return nil
	},
}

func init() {
	authCmd.Flags().StringP("host", "", "github.com", "GitHub host (default: github.com)")
	authCmd.Flags().StringP("permission", "p", "", "Optional permission scope (e.g. repo)")
	rootCmd.AddCommand(authCmd)
}