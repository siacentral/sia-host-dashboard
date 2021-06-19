package cmd

import (
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/inconshreveable/mousetrap"
)

// StartedInExplorer on windows checks if we need to spawn a new command line to prevent
// immediately closing the window. On other systems, does nothing
func StartedInExplorer() {
	if !mousetrap.StartedByExplorer() {
		return
	}

	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CmdLine: `cmd /K start dashboard.exe`,
	}

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
