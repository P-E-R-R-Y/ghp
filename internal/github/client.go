package github

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{Token: token}
}

// ======================
// Fetch user/org info
// ======================
func (c *Client) FetchUserInfo(username string) (*User, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
	return c.fetchInfo(url)
}

func (c *Client) FetchOrgInfo(org string) (*User, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s", org)
	return c.fetchInfo(url)
}

func (c *Client) fetchInfo(url string) (*User, error) {
	req, _ := http.NewRequest("GET", url, nil)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var info User
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	return &info, nil
}

// ======================
// Fetch repos
// ======================
func (c *Client) FetchUserRepos(username string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
	return c.fetchRepos(url)
}

func (c *Client) FetchOrgRepos(org string) ([]Repository, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
	return c.fetchRepos(url)
}

func (c *Client) fetchRepos(url string) ([]Repository, error) {
	req, _ := http.NewRequest("GET", url, nil)
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}
	return repos, nil
}