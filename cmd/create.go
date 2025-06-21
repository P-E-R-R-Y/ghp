package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a GitHub repository under user or organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		org, _ := cmd.Flags().GetString("org")
		private, _ := cmd.Flags().GetBool("private")
		template, _ := cmd.Flags().GetBool("template")

		if name == "" {
			return errors.New("repository --name is required")
		}

		target := "user"
		if org != "" {
			target = "organization: " + org
		}

		fmt.Printf("You are about to create a repository with the following settings:\n")
		fmt.Printf("  Name: %s\n", name)
		fmt.Printf("  Target: %s\n", target)
		fmt.Printf("  Private: %v\n", private)
		fmt.Printf("  Template: %v\n", template)
		fmt.Print("Confirm? (y/N): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		response = strings.TrimSpace(response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Aborted.")
			return nil
		}

		// Prepare request body
		body := map[string]interface{}{
			"name":    name,
			"private": private,
		}

		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}

		var endpoint string
		if org != "" {
			endpoint = fmt.Sprintf("/orgs/%s/repos", org)
		} else {
			endpoint = "/user/repos"
		}

		cmdExec := exec.Command("gh", "api", "-X", "POST", endpoint, "-H", "Accept: application/vnd.github.v3+json", "--input", "-")
		cmdExec.Stdin = bytes.NewReader(jsonBody)

		output, err := cmdExec.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create repository: %s, error: %w", string(output), err)
		}

		fmt.Println("Repository created successfully:")
		fmt.Println(string(output))
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("name", "n", "", "Name of the repository (required)")
	createCmd.Flags().StringP("org", "o", "", "Create repository in this organization")
	createCmd.Flags().BoolP("private", "p", false, "Create private repository")
	createCmd.Flags().BoolP("template", "t", false, "Mark repository as a template repository (preview API)")
	rootCmd.AddCommand(createCmd)
}