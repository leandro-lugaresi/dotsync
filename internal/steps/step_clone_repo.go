package steps

import (
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/leandro-lugaresi/dotsync/internal/vcs"
	"github.com/mitchellh/colorstring"
	"github.com/mitchellh/multistep"
)

type StepCloneRepo struct{}

func (*StepCloneRepo) Run(state multistep.StateBag) multistep.StepAction {
	repo := state.Get("repo").(string)
	path := state.Get("path").(string)
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	colorstring.Printf("[white]=> Cloning repository: %s", repo)
	s.Start()

	_, err := vcs.NewGitRepository(path, repo, os.Stdout)
	if err != nil {
		s.Stop()
		colorstring.Printf("[red]=> Error cloning the repository %s \n [white]error: %s", repo, err)
		state.Put("repo_result", "error")
		return multistep.ActionHalt
	}

	// Print a success dot
	s.Stop()
	colorstring.Print("[green]=> Clone successfully!")
	state.Put("repo_result", "clone")
	return multistep.ActionContinue
}

func (*StepCloneRepo) Cleanup(multistep.StateBag) {}
