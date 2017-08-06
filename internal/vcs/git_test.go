package vcs

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const initialCommit string = "02c39d8f899fa1fe19f0a2bf7d983ccf88314840"

var integration = flag.Bool("integration", false, "Run integration tests")

func Test_NewGitRepository(t *testing.T) {
	p := createTempDir(t)
	defer os.Remove(p)
	t.Run("ShouldCloneRepository", func(t *testing.T) {
		_, err := NewGitRepository(p, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
		check(t, "Failed To init the repository", err)
		_, err = ioutil.ReadDir(filepath.Join(p, ".git"))
		check(t, "Failed to create the directory for .git", err)
	})

	t.Run("ShouldOpenTheRepository", func(t *testing.T) {
		r, err := NewGitRepository(p, "", ioutil.Discard)
		check(t, "Failed To init the repository", err)
		remote, err := r.repo.Remote("origin")
		check(t, "Failed to get the remote origin", err)
		assert.Equal(t, "git@gitlab.com:leandro-lugaresi/git-tests.git", remote.Config().URL)
	})
}

func Test_gitPushAndPull(t *testing.T) {
	if !*integration {
		t.Skip("Skip integration test")
	}
	p1 := createTempDir(t)
	p2 := createTempDir(t)
	defer os.Remove(p1)
	defer os.Remove(p2)
	var hash [20]byte

	r1, err := NewGitRepository(p1, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	check(t, "Failed to clone the repository", err)
	r2, err := NewGitRepository(p2, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	check(t, "Failed to clone the repository", err)

	t.Run("Push commit modifications", func(t *testing.T) {
		err = ioutil.WriteFile(filepath.Join(p1, "testFoo.log"), []byte("hello world!"), 0644)
		check(t, "Failed to create the test file", err)
		hash, err = r1.Commit()
		check(t, "Failed to commit files", err)
		err = r1.Push()
		check(t, "Failed to push files", err)
	})
	t.Run("Pull modifications", func(t *testing.T) {
		err = r2.Pull()
		check(t, "Failed to push files", err)
		ref, err := r2.repo.Head()
		check(t, "Failed to get the reference for repository", err)
		assert.EqualValues(t, hash, ref.Hash())
	})

	err = r2.ForceReset()
	check(t, "Failed to reset the repository", err)
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

func createTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir(os.TempDir(), "dotfiles")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func check(t *testing.T, fail string, err error) {
	if err != nil {
		t.Fatal(fail, " - error: ", err)
	}
}
