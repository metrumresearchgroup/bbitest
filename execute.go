package babylontest

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
)

func executeCommand(ctx context.Context, command string, args... string) string{
	//Find it in path
	binary, _ := exec.LookPath(command)
	cmd := exec.CommandContext(ctx,binary, args...)
	cmd.Env = os.Environ()
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Errorf("An error occurred trying to execute model. Error details are : %s", err)

		if exitError, ok := err.(*exec.ExitError); ok {
			code := exitError.ExitCode()
			details := exitError.String()

			log.Errorf("Exit code was %d, details were %s", code, details)
			log.Errorf("output details were: %s",string(output))
		}
	}

	//log.Info(string(output))
	outputString := string(output)
	return outputString
}