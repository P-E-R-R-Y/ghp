package github

type User struct {
	Login       string `json:"login"`
	Name        string `json:"name"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
}

type Repository struct {
	ID	 int64  `json:"id"`
	Name string `json:"name"`
	HTMLURL     string `json:"html_url"`
	Description string `json:"description,omitempty"`
	Language    string `json:"language,omitempty"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
}