package job

import (
	"fmt"
	"os/exec"
	"strings"
)

func Command(c string) ([]byte, int, error) {
	cmdFull := strings.Split(c, " ")
	cmd := cmdFull[0]
	out, err := exec.Command(cmd, cmdFull[1:]...).CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return out, exitError.ExitCode(), fmt.Errorf("error in command execution: %s", err.Error())
		} else {
			return out, -1, fmt.Errorf("error in command execution: %s", err.Error())
		}
	}
	return out, 0, nil
}
