package steps

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
)

type gitRepository struct {
	output io.Writer
	repo   *git.Repository
}

func (g *gitRepository) push() error {
	return g.repo.Push(&git.PushOptions{})
}

func (g *gitRepository) pull() error {
	return g.repo.Pull(&git.PullOptions{
		Progress:          g.output,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
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
