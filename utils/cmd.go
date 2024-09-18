package utils

import (
	"os/exec"

	"github.com/oddmario/tunnel-manager/logger"
)

func Cmd(cmd string, shell, log_errors bool) ([]byte, error) {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()

		if err != nil {
			if log_errors {
				logger.Error("Command \"" + cmd + "\" failed: " + err.Error())
			}

			return nil, err
		}

		return out, nil
	} else {
		out, err := exec.Command(cmd).Output()

		if err != nil {
			if log_errors {
				logger.Error("Command \"" + cmd + "\" failed: " + err.Error())
			}

			return nil, err
		}

		return out, nil
	}
}
