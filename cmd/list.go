package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/P-E-R-R-Y/ghp/internal/github"
	"github.com/P-E-R-R-Y/ghp/internal/color"
)

var ghlistCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories for a GitHub user or organization",
	RunE: func(cmd *cobra.Command, args []string) error {

		defaultGroups := []string{"^(impobject|iobject).*", "^i.*module$", ".*module$", "^i.*", ".*"}

		user, _ := cmd.Flags().GetString("user")
		org, _ := cmd.Flags().GetString("org")
		filters, _ := cmd.Flags().GetStringSlice("filter")
		groups, _ := cmd.Flags().GetStringSlice("group")
		repeat, _ := cmd.Flags().GetBool("repeat")
		
		// Validate flags
		if user == "" && org == "" {
			return errors.New("must provide either --user or --org")
		}
		if user != "" && org != "" {
			return errors.New("only one of --user or --org can be set")
		}
	
		// Build endpoint
		var endpoint string
		if user != "" {
			endpoint = fmt.Sprintf("/users/%s/repos", user)
		} else {
			endpoint = fmt.Sprintf("/orgs/%s/repos", org)
		}
	
		// Call gh
		cmdOutput, err := exec.Command("gh", "api", endpoint, "--paginate").Output()
		if err != nil {
			return fmt.Errorf("failed to run gh: %w", err)
		}
			
		// Parse JSON
		var repos []github.Repository
		if err := json.Unmarshal(cmdOutput, &repos); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}
	
		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return nil
		}
	
		// --- Filtering ---
		var results []github.Repository
	
		if filters != nil && len(filters) > 0 {
			for _, repo := range repos {
				for i, filter := range filters {
					regex, err := regexp.Compile(filter)
					if err != nil {
						return fmt.Errorf("the filter #%d is invalid regex: %v", i, err)
					}
					if regex.MatchString(repo.Name) {
						results = append(results, repo)
						break
					}
				}
			}
		} else {
			results = repos
		}
	
		// --- Grouping & Printing ---
		printed := make(map[string]bool)
	
		if len(groups) == 0 {
			// Use defaultGroups if groups not set
			groups = defaultGroups
		}
	
		for _, group := range groups {
			re, err := regexp.Compile(group)
			if err != nil {
				return fmt.Errorf("invalid group regex %q: %v", group, err)
			}
	
			fmt.Printf("Group: %s\n", group)
			for _, repo := range results {
				if printed[repo.Name] && !repeat {
					continue
				}
				if re.MatchString(repo.Name) {
					templateStr := map[bool]string{true: color.Yellow + "template" + color.Reset, false: ""}[repo.IsTemplate]
					fmt.Printf(" - %s%s%s (%s) %s %s\n", color.Green, repo.Name, color.Blue, repo.HTMLURL, templateStr, color.Reset)
					printed[repo.Name] = true
				}
			}
			fmt.Println()
		}
	
		fmt.Println("Others:")
		for _, repo := range results {
			if !printed[repo.Name] {
				templateStr := map[bool]string{true: color.Yellow + "template" + color.Reset, false: ""}[repo.IsTemplate]
				fmt.Printf(" - %s%s%s (%s) %s %s\n", color.Green, repo.Name, color.Blue, repo.HTMLURL, templateStr, color.Reset)
			}
		}	
		return nil
	},
}

func init() {
	ghlistCmd.Flags().StringP("user", "u", "", "GitHub username")
	ghlistCmd.Flags().StringP("org", "o", "", "GitHub organisation")
	ghlistCmd.Flags().StringSliceP("filter", "f", nil, "Filter result")
	ghlistCmd.Flags().StringSliceP("group", "g", nil, "Group result")
	ghlistCmd.Flags().BoolP("repeat", "r", false, "Repeat repos in multiple groups if matched")
	rootCmd.AddCommand(ghlistCmd)
}
