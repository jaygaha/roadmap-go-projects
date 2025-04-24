package api

type GitHubActivityStrct struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"repo"`
	Payload struct {
		Action  string `json:"action"`
		Ref     string `json:"ref"`
		RefType string `json:"ref_type"`
		Commits []struct {
			Message string `json:"message"`
		} `json:"commits"`
	} `json:"payload"`
}

type GitHubFeedStrct struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Actor struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		URL       string `json:"url"`
		AvatarURL string `json:"avatar_url"`
	} `json:"actor"`
	Repo struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		AvatarURL string `json:"avatar_url"`
	} `json:"repo"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Public    bool   `json:"public"`
	Processed bool   `json:"processed"`
}
