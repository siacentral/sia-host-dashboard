package cmd

import (
	"fmt"
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
		CmdLine: fmt.Sprintf(`cmd /K start "%s"`, os.Args[0]),
	}

	if err := cmd.Run(); err != nil {
		log.Fatalln(err)
	}

	os.Exit(0)
}
