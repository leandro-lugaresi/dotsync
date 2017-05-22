// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

// A vcsCmd describes how to use a version control system
// like Mercurial, Git, or Subversion, this is based on vcsCmd for go get.
type vcsCmd struct {
	name string
	cmd  string // name of binary to invoke command

	createCmd   []string // commands to download a fresh copy of a repository
	downloadCmd []string // commands to download updates into an existing repository
	publishCmd  []string // commands to publish changes into an existing repository

	scheme  []string
	pingCmd string

	remoteRepo  func(v *vcsCmd, rootDir string) (remoteRepo string, err error)
	resolveRepo func(v *vcsCmd, rootDir, remoteRepo string) (realRepo string, err error)
}

func (cmd *vcsCmd) run() {

}
