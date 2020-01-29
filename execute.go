package babylontest

import (
	"context"
	"os"
	"os/exec"
)

func executeCommand(ctx context.Context, command string, args... string) string{
	//Find it in path
	binary, _ := exec.LookPath(command)
	cmd := exec.CommandContext(ctx,binary, args...)
	cmd.Env = os.Environ()
	//log.Infof("Command is %s", cmd.String())
	output, _ := cmd.CombinedOutput()
	//log.Info(string(output))

	return string(output)
}