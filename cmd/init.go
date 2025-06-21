package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize environment and check dependencies",
	RunE: func(cmd *cobra.Command, args []string) error{
		if !isGhInstalled() {
			switch runtime.GOOS {
			case "darwin":
				fmt.Println("🚨 GitHub CLI (gh) is not installed.")
				fmt.Println("💡 You can install it using Homebrew: `brew install gh`")
			default:
				fmt.Println("🚨 GitHub CLI (gh) is not installed.")
				fmt.Println("📦 Download it from: https://cli.github.com/")
			}
			return nil
		}

		fmt.Println("✅ GitHub CLI (gh) is installed.")
		return nil
	},
}

func isGhInstalled() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}