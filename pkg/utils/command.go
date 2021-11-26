package utils

import (
	"fmt"
	"os/exec"
)

func RunCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("run command error, err:%s, cmd:%s, args: %s\n", err, name, args)
		fmt.Println(string(out))
		return nil, err
	}

	return out, nil
}
