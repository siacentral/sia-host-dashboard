package cmd

// StartedInExplorer on windows checks if we need to spawn a new command line to prevent
// immediately closing the window. On other systems, does nothing
func StartedInExplorer() {
	return
}
