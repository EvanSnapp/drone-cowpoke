package types

type CatalogInfo struct {
	CatalogRepo string `yaml:"catalogRepo"`
	Version     int    `yaml:"version"`
	Branch      string `yaml:"branch"`
}
