package steps

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type gitRepository struct {
	output io.Writer
	repo   *git.Repository
}

func (g *gitRepository) push() error {
	return g.repo.Push(&git.PushOptions{})
}

func (g *gitRepository) pull() error {
	tree, err := g.repo.Worktree()
	if err != nil {
		return err
	}

	return tree.Pull(&git.PullOptions{
		Progress:          g.output,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
}

func (g *gitRepository) commit() error {
	w, err := g.repo.Worktree()
	if err != nil {
		return err
	}
	s, err := w.Status()
	if err != nil {
		return err
	}

	if !s.IsClean() {
		for path, status := range s {
			switch status.Worktree {
			case git.Unmodified:
				continue
			case git.Added, git.Modified:
				_, err = w.Add(path)
				if err != nil {
					return err
				}
			case git.Deleted:
				_, err = w.Remove(path)
				if err != nil {
					return err
				}
			case git.Renamed:
			case git.Copied:
			case git.Untracked:
			case git.UpdatedButUnmerged:
			}

			if err != nil {
				return err
			}
		}

	}

	_, err = w.Commit("TODO MSG", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Dotsync Daemon",
			Email: "leandrolugaresi92@gmail.com",
			When:  time.Now(),
		},
	})
	return err
}

func initNewRepository(path, name string, output io.Writer) (*gitRepository, error) {
	url, err := parseRepositoryName(name)
	if err != nil {
		return nil, err
	}
	r, err := git.PlainInit(path, false)
	if err != nil {
		return nil, err
	}
	_, err = r.CreateRemote(&config.RemoteConfig{
		URL:  url,
		Name: "origin",
	})
	if err != nil {
		return nil, err
	}
	return &gitRepository{
		output: output,
		repo:   r,
	}, nil
}

func openRepository(path string, output io.Writer) (*gitRepository, error) {
	r, err := git.PlainOpen(path)
	return &gitRepository{
		output: output,
		repo:   r,
	}, err
}

func cloneRepository(path, name string, output io.Writer) (*gitRepository, error) {
	url, err := parseRepositoryName(name)
	if err != nil {
		return nil, err
	}
	r, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Progress:          output,
	})
	return &gitRepository{
		output: output,
		repo:   r,
	}, err
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
