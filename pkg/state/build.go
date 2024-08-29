package state

import (
	"path"

	"github.com/moby/buildkit/client/llb"
)

const keyserver = "keys.openpgp.org"

func BuildBase(base string) *State {
	s := From(base)

	s.Install(
		"autoconf",
		"build-essential",
		"curl",
		"git",
		"gnupg2",
		"libssl-dev",
		"libyaml-dev",
		"libreadline6-dev",
		"libncurses5-dev",
		"libffi-dev",
		"libgdbm-dev",
		"zlib1g-dev",
	)

	return s
}

func (s *State) GPGVerify(file string, verifyFile string, keys []string) {
	// TODO: right now, there is no dependency of s.current on execState, so this
	// will be completely ignored :/
	filePath := path.Join("/dest", file)
	verifyFilePath := path.Join("/dest", verifyFile)
	execState := llb.Image(s.baseImage).Run(recvKeysCmd(keys))
	execState = execState.Run(verifyCmd(filePath, verifyFilePath))
	execState.AddMount("/src", s.current)
}

func recvKeysCmd(keys []string) llb.RunOption {
	cmd := []string{
		"/usr/bin/gkg",
		"--batch",
		"--keyserver",
		keyserver,
		"--recv-keys",
	}

	return llb.Args(append(cmd, keys...))
}

func verifyCmd(file string, verifyFile string) llb.RunOption {
	return llb.Args([]string{"/usr/bin/gpg", "--batch", "--verify", verifyFile, file})
}
