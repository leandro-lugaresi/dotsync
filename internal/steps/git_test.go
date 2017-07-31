package steps

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func Test_gitSimpleWorkflow(t *testing.T) {
	p1 := path()
	p2 := path()

	r1, err := cloneRepository(p1, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	CheckIfError(t, "Failed to clone the repository", err)
	_, err = ioutil.ReadDir(filepath.Join(p1, ".git"))
	CheckIfError(t, "Failed to create a .git directory", err)

	_, err = cloneRepository(p2, "gitlab.com/leandro-lugaresi/git-tests", ioutil.Discard)
	CheckIfError(t, "Failed to clone the repository", err)
	_, err = ioutil.ReadDir(filepath.Join(p2, ".git"))
	CheckIfError(t, "Failed to create a .git directory", err)

	err = ioutil.WriteFile(filepath.Join(p1, "testFoo.log"), []byte("hello world!"), 0644)
	CheckIfError(t, "Failed to create the test file", err)
	err = r1.add("testFoo.log")
	CheckIfError(t, "Failed to commit files", err)

	remote, err := r1.repo.Remote("origin")
	CheckIfError(t, "Failed to get the remote origin", err)
	assert.Equal(t, "git@gitlab.com:leandro-lugaresi/git-tests.git", remote.Config().URL)

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
