package scripts

import (
	"alshx/src/platform/files"
	"alshx/src/platform/logging"
	"os"
	"os/exec"
	"path/filepath"
)

func execNode(
	logger *logging.Logger,
	cmd *exec.Cmd,
	config *Meta,
	args []string,
	folderPath string,
) {
	if !hasCommand("node") || !hasCommand("npx") {
		logger.Log("Node is not installed")
		return
	}
	if !hasCommand("yarn") {
		logger.Log("Yarn is not installed")
		return
	}
	installNodeModules(logger, folderPath)
	cmdPath := []string{"node", config.Entrypoint}
	cmdPath = append(cmdPath, config.Args...)
	cmdPath = append(cmdPath, args...)
	logger.Info("Command:", cmdPath)
	cmd = exec.Command(cmdPath[0], cmdPath[1:]...)
}

func installNodeModules(
	logger *logging.Logger,
	folderPath string,
) {
	if files.NotExists(filepath.Join(folderPath, "node_modules")) {
		logger.Log("Installing Node Modules")
		cmd := exec.Command("yarn")
		cmd.Dir = folderPath
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
