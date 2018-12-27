package repo

import (
	"fmt"
	"github.com/pkg/errors"
	"path"
	"sync"

	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"

	"github.com/GitbookIO/monorepo/monofile"
)

type Repo struct {
	Directory    string
	CachePath    string
	Monofile     *monofile.Monofile
	PathMonofile string
	Lockfile     *monofile.LockedMonofile
	PathLockfile string
}

func Open(dir string) (*Repo, error) {
	mpath := path.Join(dir, "monorepo.yml")
	lpath := path.Join(dir, "monorepo.lock")

	// Read monofile.yml
	m, err := monofile.ReadFile(mpath)
	if err != nil {
		return nil, err
	}

	// Lock (optional)
	var l *monofile.LockedMonofile
	if lf, err := monofile.ReadLockFile(lpath); err == nil {
		l = lf
	}

	return &Repo{
		Directory:    dir,
		CachePath:    path.Join(dir, ".cache"),
		Monofile:     m,
		PathMonofile: mpath,
		Lockfile:     l,
		PathLockfile: lpath,
	}, nil
}

func (r *Repo) List() error {
	if r.Lockfile != nil {
		r.listLocked()
	} else {
		r.listSimple()
	}

	return nil
}

func (r *Repo) listSimple() error {
	for _, sub := range r.Monofile.Repos {
		fmt.Printf("%s (%s) - (%s)\n", sub.Path, sub.URL, sub.Ref)
	}
	return nil
}

func (r *Repo) listLocked() error {
	for _, sub := range r.Lockfile.Repos {
		fmt.Printf("%s (%s) - (%s)[%s]\n", sub.Path, sub.URL, sub.Ref, sub.SHA)
	}
	return nil
}

// Pushes all repos
func (r *Repo) Push(force bool) error {
	return nil
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
func (r *Repo) Pull(force bool) error {
	wg := sync.WaitGroup{}
	for _, sub := range r.Monofile.Repos {
		wg.Add(1)
		go func(sub monofile.Repo) {
			defer wg.Done()
			err := r.pullsub(sub, force)
			if err != nil {
				fmt.Printf("Error pulling '%s': %s", sub.URL, err)
			}
		}(sub)
	}
	wg.Wait()
	return nil
}

// PullSub updates a specific subrepo
func (r *Repo) PullSub(pathOrURL string, force bool) error {
	sub := r.lookupRepo(pathOrURL)
	if sub == nil {
		return fmt.Errorf("PullSub('%s'): no such subrepo", pathOrURL)
	}
	return r.pullsub(*sub, force)
}

func (r *Repo) pullsub(sub monofile.Repo, force bool) error {
	// Fetch changes
	if err := r.cacheUpdate(sub.URL, sub.Path, sub.Ref); err != nil {
		return errors.Wrap(err, "pullsub/cacheUpdate")
	}

	// Get SHA
	sha, err := r.lookupSHA(sub.Path, sub.Ref)
	if err != nil {
		return err
	}

	// Gen locked info
	locked := monofile.LockedRepo{
		Repo: sub,
		SHA:  sha,
	}

	// Do checkout
	if err := r.checkoutsub(locked); err != nil {
		return errors.Wrap(err, "pullsub/checkout")
	}

	return nil
}

// Save updates the monofiles
func (r *Repo) Save() error {
	return nil
}

// Updates a specific repo in the lockfile
func (r *Repo) updateLockfile(sub monofile.Repo, sha string) error {
	return nil
}

func (r *Repo) lookupSHA(subkey, ref string) (string, error) {
	// Open repo
	gitrepo, err := r.gitrepo(subkey)
	if err != nil {
		return "", errors.Wrapf(err, "lookupSHA/open('%s')", subkey)
	}

	// Lookup hash
	hash, err := gitrepo.ResolveRevision(plumbing.Revision("refs/heads/" + ref))
	if err != nil {
		return "", errors.Wrapf(err, "lookupSHA/resolve('%s', '%s')", subkey, ref)
	}

	return hash.String(), nil
}

// Add with clone the subrepo and add to the monofile
func (r *Repo) Add(url, subkey, ref string) error {
	// Do clone
	if err := r.cacheUpdate(url, subkey, ref); err != nil {
		return errors.Wrap(err, "add/clone")
	}

	// Get SHA
	sha, err := r.lookupSHA(subkey, ref)
	if err != nil {
		return err
	}

	// Gen locked info
	locked := monofile.LockedRepo{
		Repo: monofile.Repo{
			URL:  url,
			Path: subkey,
			Ref:  ref,
		},
		SHA: sha,
	}

	// Do checkout
	if err := r.checkoutsub(locked); err != nil {
		return errors.Wrap(err, "add/checkout")
	}

	// Add to lockfile
	err = r.updateLockfile(locked.Repo, locked.SHA)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) checkoutsub(sub monofile.LockedRepo) error {
	// Open repo
	gitrepo, err := r.gitrepo(sub.Path)

	// Worktree
	wt, err := gitrepo.Worktree()
	if err != nil {
		return err
	}
	// Do checkout
	err = wt.Reset(&git.ResetOptions{
		Commit: plumbing.NewHash(sub.SHA),
		Mode:   git.MixedReset,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) gitrepo(subkey string) (*git.Repository, error) {
	gitdir := r.cachedir(subkey)
	cachefs, err := filesystem.NewStorage(osfs.New(gitdir))
	if err != nil {
		return nil, err
	}
	// Path of working dir to checkout to
	wtdir := r.workingdir(subkey)

	// Open repo
	gitrepo, err := git.Open(cachefs, osfs.New(wtdir))
	if err != nil {
		return nil, err
	}

	return gitrepo, nil
}

// pullOrClone fetches a repo into the local cache
func (r *Repo) cacheUpdate(url, subkey, ref string) error {
	// Try to open and fetch
	err := r.cacheFetch(url, subkey, ref)
	if err == git.ErrRepositoryNotExists {
		// Fallback to cloning
		err = nil
		if err := r.cacheClone(url, subkey, ref); err != nil {
			return errors.Wrap(err, "cacheUpdate/clone")
		}
	}
	// Check error
	if err != nil {
		return errors.Wrap(err, "cacheUpdate/fetch")
	}

	return nil
}

func (r *Repo) cacheFetch(url, subkey, ref string) error {
	gitrepo, err := git.PlainOpen(r.cachedir(subkey))
	if err != nil {
		return err
	}

	// Do fetch
	err = gitrepo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			"+refs/heads/*:refs/remotes/origin/*",
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "cacheFetch/fetch")
	}

	return nil
}

func (r *Repo) cacheClone(url, subkey, ref string) error {
	isBare := true
	_, err := git.PlainClone(
		r.cachedir(subkey),
		isBare,
		&git.CloneOptions{
			URL: url,
		},
	)
	return err
}

func (r *Repo) cachedir(subkey string) string {
	return path.Join(r.CachePath, subkey)
}

func (r *Repo) workingdir(subkey string) string {
	return path.Join(r.Directory, subkey)
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
	for _, r := range r.Lockfile.Repos {
		if repoMatch(r.Repo, pathOrURL) {
			return &r
		}
	}
	return nil
}

func repoMatch(r monofile.Repo, pathOrURL string) bool {
	return r.Path == pathOrURL || r.URL == pathOrURL
}
