package cmd

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/P-E-R-R-Y/gitperry/internal/github"
	"github.com/P-E-R-R-Y/gitperry/internal/color"
	"github.com/P-E-R-R-Y/gitperry/config"

	"github.com/spf13/cobra"

)

var filters []string
var groups []string
var defaultGroups = []string{"^i.*module$", ".*module$", "^i.*", ".*"}
var repeat bool


var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List repositories for a GitHub user or organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags: either org or user must be set, but not both
		if (user == "") {
			return errors.New("must specify --org or --user")
		}

		token := config.GetToken()
		client := github.NewClient(token)

		repos, err := client.FetchUserRepos(user)
		
		if err != nil {
			return fmt.Errorf("error fetching repositories: %w", err)
		}

		if len(repos) == 0 {
			fmt.Println("No repositories found.")
			return nil
		}

		//Filter
		var results []github.Repository

		if (filters != nil) {
			for _, repo := range repos {
				for i, filter := range filters {
					regex, err := regexp.Compile(filter)
					if err != nil {
						return errors.New(fmt.Sprintf("the filter nb: %i, failed to give a regex", i))
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

		//groups
		printed := make(map[string]bool) // track printed repo names

		if len(groups) == 0 {
			// --group was passed but empty (i.e. "--group" with no value)
			groups = defaultGroups
		}
		for _, group := range groups {
			re := regexp.MustCompile(group)
			fmt.Printf("Group: %s\n", group)
			for _, repo := range results {
				if printed[repo.Name] && !repeat {
					continue // skip if already printed
				}
				if re.MatchString(repo.Name) {
					fmt.Printf(" - %s%s%s (%s)%s\n", color.Green, repo.Name, color.Blue, repo.HTMLURL, color.Reset)
					printed[repo.Name] = true
				}
			}
			fmt.Println()
		}

		fmt.Printf("Others:\n")
		for _, repo := range results {
			if !printed[repo.Name] {
				fmt.Printf(" - %s (%s)\n", repo.Name, repo.HTMLURL)
			}
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringVarP(&user, "user", "u", "", "GitHub username")
	listCmd.Flags().StringSliceVarP(&filters, "filter", "f", nil, "Filter result")
	listCmd.Flags().StringSliceVarP(&groups, "group", "g", nil, "Group result")
	listCmd.Flags().BoolVarP(&repeat, "repeat", "r", false, "Repeat value if match with multiple groups")
	rootCmd.AddCommand(listCmd)
}
