package types

// contains plugin arguments
type Plugin struct {
	CatalogRepo        string `json:"catalog_repo"`
	GitHubToken        string `json:"github_token"`
	CowpokeURL         string `json:"cowpoke_url"`
	RancherCatalogName string `json:"rancher_catalog_name"`
	BearerToken        string `json:"bearer_token"`
}
