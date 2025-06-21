package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"github.com/P-E-R-R-Y/ghp/internal/color"
)

var deleteRepoCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a GitHub repository (requires confirmation)",
	RunE: func(cmd *cobra.Command, args []string) error {
		//user
		user, _ := cmd.Flags().GetString("org")
		if user == "" {
			// Get current authenticated user
			output, err := exec.Command("gh", "api", "/user", "--jq", ".login").Output()
			if err != nil {
				return fmt.Errorf("failed to get current user: %w", err)
			}
			user = strings.TrimSpace(string(output))
		}
		//repo
		repo, _ := cmd.Flags().GetString("repo")
		if repo == "" {
			return errors.New("please specify --repo in format owner/name")
		}

		fmt.Printf("WARNING: You are about to delete the repository %s'%s/%s'%s. This action is irreversible.\n", color.Red, user, repo, color.Reset)
		fmt.Printf("To confirm, please type the repository name %s'%s/%s'%s: ", color.Red, user, repo, color.Reset)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read confirmation: %w", err)
		}

		input = strings.TrimSpace(input)
		if strings.ToLower(input) != strings.ToLower(user + "/" + repo) {
			fmt.Println("Confirmation did not match repository name. Aborting.")
			return nil
		}

		fmt.Println("deleting repository...")

		// Run gh API delete command
		cmdExec := exec.Command("gh", "api", "-X", "DELETE", fmt.Sprintf("/repos/%s/%s", user, repo))
		output, err := cmdExec.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to delete repository: %s, error: %w", string(output), err)
		}

		fmt.Printf("Repository '%s' deleted successfully.\n", repo)
		return nil
	},
}

func init() {
	deleteRepoCmd.Flags().StringP("org", "o", "", "Organisation")
	deleteRepoCmd.Flags().StringP("repo", "r", "", "Repository to delete")
	
	rootCmd.AddCommand(deleteRepoCmd)
}