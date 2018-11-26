package repo

import (
	"fmt"
	"path"

	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	"github.com/GitbookIO/monorepo/monofile"
)

type Repo struct {
	Directory      string
	Monofile       *monofile.Monofile
	LockedMonofile *monofile.LockedMonofile
}

func Open(dir string) (*Repo, error) {

}

// Pushes all repos
func (r *Repo) Push(force bool) error {

}

func (r *Repo) PushSub(pathOrURL string, force bool) error {
	subrepo := r.lookupRepo(pathOrURL)
	if subrepo == nil {
		return fmt.Errorf("PushSub('%s'): no such subrepo", pathOrURL)
	}
	return r.pushsub(*subrepo, force)
}

func (r *Repo) pushsub(sub monofile.Repo, force bool) error {
	return nil
}

// Pull updates all repos
func (r *Repo) Pull(force bool) error {}

// PullSub updates a specific subrepo
func (r *Repo) PullSub(pathOrURL string, force bool) error {

}

// Updates a specific repo in the lockfile
func (r *Repo) updateLockfile(sub monofile.Repo, sha string) error {
	return nil
}

// Add with clone the subrepo and add to the monofile
func (r *Repo) Add(url, subdir, ref string) error {
	wtdir := path.Join(r.Directory, ".cache", subdir)
	gitdir := path.Join(r.Directory, ".cache", subdir, ".git")
	// Git storage
	storage, err := filesystem.NewStorage(osfs.New(gitdir))
	if err != nil {
		return err
	}
	// Do clone
	gitrepo, err := git.Clone(
		storage,
		osfs.New(wtdir),
		&git.CloneOptions{
			URL: url,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) lookupRepo(pathOrURL string) *monofile.Repo {
	for _, r := range r.Monofile.Repos {
		if repoMatch(r, pathOrURL) {
			return &r
		}
	}
	return nil
}

func (r *Repo) lookupLockedRepo(pathOrURL string) *monofile.LockedRepo {
	for _, r := range r.LockedMonofile.Repos {
		if repoMatch(r.Repo, pathOrURL) {
			return &r
		}
	}
	return nil
}

func repoMatch(r monofile.Repo, pathOrURL string) bool {
	return r.Path == pathOrURL || r.URL == pathOrURL
}
