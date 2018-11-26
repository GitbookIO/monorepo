package monofile

// Monofile is a monorepo config file
type Monofile struct {
	Repos []Repo `json:"repos" yaml:"repos"`
}

// Repo desribes a repo synced with monorepo
type Repo struct {
	Path string `json:"path" yaml:"path"`
	URL  string `json:"url" yaml:"url"`
	Ref  string `json:"ref" yaml:"ref"`
}

// LockedMonofile is a lockfile, locking repos to SHAs
type LockedMonofile struct {
	Repos []LockedRepo `json:"repos" yaml:"repos"`
}

// LockedRepo is a repo locked to a SHA
type LockedRepo struct {
	Repo
	SHA string `json:"sha" yaml:"sha"`
}
