package cmd

import (
	"errors"
	"fmt"

	"github.com/P-E-R-R-Y/gitperry/internal/github"

	"github.com/P-E-R-R-Y/gitperry/config"

	"github.com/spf13/cobra"
)

var user string

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get info about a GitHub user or organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate flags: exactly one of --user or --org must be set
		if (user == "") {
			return errors.New("must specify --user example")
		}

		token := config.GetToken()
		client := github.NewClient(token)

		if client == nil {
			return errors.New("client creation failed")
		}

		userInfo, err := client.FetchUserInfo(user)
		if err != nil {
			return fmt.Errorf("error fetching user info: %w", err)
		}
		printUserInfo(*userInfo)
		return nil
	},
}

func init() {
	infoCmd.Flags().StringVarP(&user, "user", "u", "", "GitHub username or organisation (mandatory)")
	rootCmd.AddCommand(infoCmd)
}

func printUserInfo(u github.User) {
	fmt.Println("User Info:")
	fmt.Printf("Login: %s\nName: %s\nBio: %s\nPublic Repos: %d\nFollowers: %d\n",
		u.Login, u.Name, u.Bio, u.PublicRepos, u.Followers)
}
