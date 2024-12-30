package state

import (
	"path"

	"github.com/moby/buildkit/client/llb"
)

func (s *State) Env(key string, value string) {
	s.current = s.current.AddEnv(key, value)
}

func (s *State) CreateDirectory(path string) {
	s.Exec("/bin/mkdir", "-p", path)
}

func (s *State) SymLink(src string, target string) {
	s.Exec("/bin/ln", "-s", src, target)
}

func (s *State) Install(pkgs ...string) {
	execState := s.current.Run(updateCmd()) // TODO: only update if necessary (how?)
	execState = execState.Run(installCmd(pkgs))
	s.current = execState.Root()
}

func (s *State) Copy(src *State, srcPath string, destPath string) {
	destDir := path.Join("/dest", path.Dir(destPath))
	execState := llb.Image(s.baseImage).Run(mkdirCmd(destDir))
	s.current = execState.AddMount("/dest", s.current)

	execState = s.current.Run(cpCmd(path.Join("/src", srcPath), path.Join("/dest", destPath)))
	execState.AddMount("/src", src.current, llb.Readonly)
	s.current = execState.AddMount("/dest", s.current)
}

func (s *State) CopyDir(src *State, srcPath string, destPath string) {
	execState := llb.Image(s.baseImage).Run(llb.Args([]string{"/bin/mkdir", "-p", path.Join("/dest", path.Dir(destPath))}))
	s.current = execState.AddMount("/dest", s.current)

	execState = s.current.Run(llb.Args([]string{"/bin/cp", "-a", "-R", path.Join("/src", srcPath) + "/.", "-t", path.Join("/dest", destPath)}))
	execState.AddMount("/src", src.current, llb.Readonly)
	s.current = execState.AddMount("/dest", s.current)
}

func (s *State) Download(url string, destPath string) {
	destDir := path.Join("/dest", path.Dir(destPath))
	execState := llb.Image(s.baseImage).Run(mkdirCmd(destDir))
	execState = execState.Run(updateCmd())
	execState = execState.Run(installCmd([]string{"apt-transport-https", "curl", "ca-certificates"}))
	execState = execState.Run(curlCmd(url, path.Join("/dest", destPath)))

	s.current = execState.AddMount("/dest", s.current)
}

func updateCmd() llb.RunOption {
	return llb.Args([]string{"/usr/bin/apt-get", "update"})
}

func installCmd(pkgs []string) llb.RunOption {
	cmd := []string{
		"/usr/bin/apt-get",
		"install",
		"-y",
		"--no-install-recommends",
		"--no-install-suggests",
	}

	return llb.Args(append(cmd, pkgs...))
}

func mkdirCmd(dir string) llb.RunOption {
	return llb.Args([]string{"/bin/mkdir", "-p", dir})
}

func cpCmd(src string, target string) llb.RunOption {
	return llb.Args([]string{"/bin/cp", "-a", "-T", src, target})
}

func curlCmd(url string, destPath string) llb.RunOption {
	return llb.Args([]string{"/usr/bin/curl", "-L", "-o", destPath, url})
}
