package steps

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/multistep"
)

type CopyConfigs struct{}

func (s *CopyConfigs) Run(state multistep.StateBag) multistep.StepAction {
	// path := state.Get("path").(string)

	return multistep.ActionContinue
}

func (s *CopyConfigs) Cleanup(multistep.StateBag) {}

// CopyFile copies the contents from src to dst atomically.
// If dst does not exist, CopyFile creates it with permissions perm.
// If the copy fails, CopyFile aborts and dst is preserved.
// Based on kelseyhightower work on https://github.com/golang/go/issues/8868.
func CopyFile(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	err = os.Rename(tmp.Name(), dst)
	if err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}
