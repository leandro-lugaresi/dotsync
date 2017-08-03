package vcs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/src-d/go-billy.v3/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type gitRepository struct {
	output io.Writer
	path   string
	repo   *git.Repository
}

func (g *gitRepository) Clone(name string) error {
	url, err := parseRepositoryName(name)
	if err != nil {
		return err
	}
	g.repo, err = git.PlainClone(g.path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          g.output,
	})
	return err
}

func (g *gitRepository) Commit() ([20]byte, error) {
	w, err := g.repo.Worktree()
	if err != nil {
		return plumbing.ZeroHash, err
	}
	s, err := w.Status()
	if err != nil {
		return plumbing.ZeroHash, err
	}

	if !s.IsClean() {
		for path, status := range s {
			switch status.Worktree {
			case git.Unmodified:
				continue
			case git.Added, git.Modified, git.Renamed, git.Copied, git.Untracked:
				_, err = w.Add(path)
				if err != nil {
					return plumbing.ZeroHash, err
				}
			case git.Deleted:
				_, err = w.Remove(path)
				if err != nil {
					return plumbing.ZeroHash, err
				}
			case git.UpdatedButUnmerged:
				continue //?
			}
		}
	}

	hash, err := w.Commit("Automatic commit for:\n\n"+s.String(), &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Dotsync Daemon",
			Email: "leandrolugaresi92@gmail.com",
			When:  time.Now(),
		},
	})
	return hash, err
}

func (g *gitRepository) Push() error {
	return g.repo.Push(&git.PushOptions{})
}

func (g *gitRepository) Pull() error {
	tree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	return tree.Pull(&git.PullOptions{
		Progress:          g.output,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}

// ForceReset is only used by tests to reset the test repository to the initial commit.
func (g *gitRepository) ForceReset() error {
	h := plumbing.NewHash(initialCommit)
	tree, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	err = tree.Reset(&git.ResetOptions{
		Commit: h,
		Mode:   git.HardReset,
	})
	if err != nil {
		return err
	}
	err = g.repo.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec("+refs/heads/master:refs/heads/master"),
		},
	})
	return err
}

// NewGitRepository return a new Repository used to do the principal operations to:
// init,clone, commit, push and pull from git repositories.
func NewGitRepository(path, name string, output io.Writer) (*gitRepository, error) {
	fs := osfs.New(path)
	_, err := fs.Stat(".git")

	g := &gitRepository{
		output: output,
		path:   path,
	}

	if os.IsNotExist(err) {
		err = g.Clone(name)
	} else {
		g.repo, err = git.PlainOpen(path)
	}
	return g, err
}

func parseRepositoryName(repo string) (string, error) {
	if strings.HasPrefix(repo, "https://") || strings.HasPrefix(repo, "git@") {
		return repo, nil
	}
	root := "github.com"
	path := repo
	for _, host := range []string{"github.com", "gitlab.com", "bitbucket.org"} {
		if strings.HasPrefix(repo, host) {
			root = host
			path = strings.Replace(repo, host+"/", "", 1)
			break
		}
	}
	if len(strings.Split(path, "/")) != 2 {
		return "", errors.New("Invalid git path, expecting foo/bar and receive " + path)
	}

	return fmt.Sprintf("git@%s:%s.git", root, path), nil
}
