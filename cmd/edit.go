package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit",
	RunE: func(cmd *cobra.Command, args []string) error {
		user, _ := cmd.Flags().GetString("org")
		name, _ := cmd.Flags().GetString("name")

		url := fmt.Sprintf("/repos/%s/%s", user, name)

		rename, _ := cmd.Flags().GetString("rename")
		template, _ := cmd.Flags().GetString("template")
		visibility, _ := cmd.Flags().GetString("visibility")

		//get user name if no user/org given
		if user == "" {
			// Get current authenticated user
			output, err := exec.Command("gh", "api", "/user", "--jq", ".login").Output()
			if err != nil {
				return fmt.Errorf("failed to get current user: %w", err)
			}
			user = strings.TrimSpace(string(output))
		}
		//control that a repo name has been given
		if name == "" {
			return errors.New("please specify --name repository_name")
		}
		//throw if no action given
		if rename == "" && template == "" && visibility == "" {
			return errors.New("please specify an action")
		}

		cmdArgs := []string{
			"api", "-X", "PATCH", url,
			"-H", "Accept:application/vnd.github+json",
		}

		//visibility
		if visibility != "" {
			if visibility != "on" && visibility != "off" {
				return errors.New("--visibility must be 'on' or 'off'.")
			}
			isPrivate := map[string]string{"on":"false", "off":"true"}[visibility]
			cmdArgs = append(cmdArgs, "-f", "private=" + isPrivate)
		}

		if template != "" {
			if template != "on" && template != "off" {
				return errors.New("--template must be 'on' or 'off'.")
			}
			isTemplate := map[string]string{"on":"true", "off":"false"}[template]
			cmdArgs = append(cmdArgs, "-f", "is_template=" + isTemplate)
		}

		if rename != "" {
			cmdArgs = append(cmdArgs, "-f", "name=" + rename)
		}

		fmt.Println(cmdArgs)
		cmdExec := exec.Command("gh", cmdArgs...)

		output, err := cmdExec.CombinedOutput()
		fmt.Println(string(output))
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	editCmd.Flags().StringP("org", "o", "", "Organisation (optional)")
	editCmd.Flags().StringP("name", "n", "", "Name of the repo to edit")
	editCmd.Flags().StringP("rename", "r", "", "Rename the repo (action)")
	editCmd.Flags().StringP("template", "t", "", "Make it a template 'on' or 'off' (action)")
	editCmd.Flags().StringP("visibility", "v", "", "Make it private 'on' or 'off' (action)")
	
	rootCmd.AddCommand(editCmd)
}