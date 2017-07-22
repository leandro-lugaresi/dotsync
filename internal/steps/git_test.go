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
	if err != nil {
		t.Error(err)
	}
	_, err = ioutil.ReadDir(filepath.Join(p, ".git"))
	if err != nil {
		t.Error(err)
	}
	g, err := openRepository(p, ioutil.Discard)
	if err != nil {
		t.Error(err)
	}
	remote, err := g.repo.Remote("origin")
	if err != nil {
		t.Error(err)
	}
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
