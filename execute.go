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

	output, _ := cmd.CombinedOutput()
	log.Info(string(output))

	return string(output)
}