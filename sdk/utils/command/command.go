package command

import (
	"io/ioutil"
	"os/exec"
)

func RunCommand(s string) (int, string, string, error) {
	cmd := exec.Command("/bin/bash", "-c", "export LANG=en_US.utf8 ; "+s)

	StdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", "", err
	}

	StderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return 0, "", "", err
	}

	exitCode := 0
	err = cmd.Start()
	if err != nil {
		return 0, "", "", err
	}

	b1, err := ioutil.ReadAll(StdoutPipe)
	if err != nil {
		return 0, "", "", err
	}
	stdout := string(b1)

	b2, err := ioutil.ReadAll(StderrPipe)
	if err != nil {
		return 0, "", "", err
	}
	stderr := string(b2)

	err = cmd.Wait()
	if err != nil {
		e, ok := err.(*exec.ExitError)
		if !ok {
			return 0, "", "", err
		}
		exitCode = e.ExitCode()
	}

	return exitCode, stdout, stderr, nil
}
