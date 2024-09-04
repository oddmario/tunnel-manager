package utils

import "os/exec"

func Cmd(cmd string, shell bool) ([]byte, error) {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			return nil, err
		}

		return out, nil
	} else {
		out, err := exec.Command(cmd).Output()

		if err != nil {
			return nil, err
		}

		return out, nil
	}
}
