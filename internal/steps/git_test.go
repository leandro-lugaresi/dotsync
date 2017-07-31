package steps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const initialCommit string = "02c39d8f899fa1fe19f0a2bf7d983ccf88314840"

func Test_gitInitAndOpen(t *testing.T) {
	p := path()
	_, err := initNewRepository(p, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	CheckIfError(t, "Failed To init the repository", err)
	_, err = ioutil.ReadDir(filepath.Join(p, ".git"))
	CheckIfError(t, "Failed to create the directory for .git", err)
	g, err := openRepository(p, ioutil.Discard)
	CheckIfError(t, "Failed to open the repository", err)
	remote, err := g.repo.Remote("origin")
	CheckIfError(t, "Failed to get the remote origin", err)
	assert.Equal(t, "git@gitlab.com:leandro-lugaresi/git-tests.git", remote.Config().URL)
}

func Test_gitPushAndPull(t *testing.T) {
	p1 := path()
	p2 := path()

	r1, err := cloneRepository(p1, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	CheckIfError(t, "Failed to clone the repository", err)
	_, err = ioutil.ReadDir(filepath.Join(p1, ".git"))
	CheckIfError(t, "Failed to create a .git directory", err)
	remote, err := r1.repo.Remote("origin")
	CheckIfError(t, "Failed to get the remote origin", err)
	assert.Equal(t, "git@gitlab.com:leandro-lugaresi/git-tests.git", remote.Config().URL)

	r2, err := cloneRepository(p2, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	CheckIfError(t, "Failed to clone the repository", err)
	_, err = ioutil.ReadDir(filepath.Join(p2, ".git"))
	CheckIfError(t, "Failed to create a .git directory", err)

	err = ioutil.WriteFile(filepath.Join(p1, "testFoo.log"), []byte("hello world!"), 0644)
	CheckIfError(t, "Failed to create the test file", err)
	hash, err := r1.commit()
	CheckIfError(t, "Failed to commit files", err)
	err = r1.push()
	CheckIfError(t, "Failed to push files", err)

	err = r2.pull()
	CheckIfError(t, "Failed to push files", err)
	ref, err := r2.repo.Head()
	CheckIfError(t, "Failed to get the reference for repository", err)
	assert.Equal(t, hash, ref.Hash())
	err = forceReset(r2)
	CheckIfError(t, "Failed to reset the repository", err)
}

func Test_parseRepositoryName(t *testing.T) {
	tests := []struct {
		repo    string
		want    string
		wantErr bool
	}{
		{"leandro-lugaresi/dotfiles", "git@github.com:leandro-lugaresi/dotfiles.git", false},
		{"github.com/leandro-lugaresi/dotfiles", "git@github.com:leandro-lugaresi/dotfiles.git", false},
		{"gitlab.com/leandro-lugaresi/dotfiles", "git@gitlab.com:leandro-lugaresi/dotfiles.git", false},
		{"git@gitlab.com:leandro-lugaresi/dotfiles.git", "git@gitlab.com:leandro-lugaresi/dotfiles.git", false},
		{"https://gitlab.com/leandro-lugaresi/dotfiles.git", "https://gitlab.com/leandro-lugaresi/dotfiles.git", false},
		{"fooo-s/bar/baz", "", true},
		{"", "", true},
		{"fooo", "", true},
		{"fooo/baz/", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.repo, func(t *testing.T) {
			got, err := parseRepositoryName(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRepositoryName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseRepositoryName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func path() string {
	dir, err := ioutil.TempDir(os.TempDir(), "dotfiles")
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func CheckIfError(t *testing.T, fail string, err error) {
	if err != nil {
		t.Fatal(fail, " - error: ", err)
	}
}

func forceReset(g *gitRepository) error {
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
