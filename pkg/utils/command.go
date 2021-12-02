package utils

import (
  "bytes"
  "os/exec"
)

//func RunCommand(name string, args ...string) ([]byte, error) {
//	cmd := exec.Command(name, args...)
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		fmt.Printf("run command error, err:%s, cmd:%s, args: %s\n", err, name, args)
//		fmt.Println(string(out))
//		return nil, err
//	}
//
//	return out, nil
//}

func RunCommand(s string) (string, error) {

  cmd := exec.Command("/bin/bash", "-c", s)
  var out bytes.Buffer
  cmd.Stdout = &out

  err := cmd.Run()

  return out.String(), err
}
