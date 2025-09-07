// internal/cli/execute.go
package cli

import (
	"os"

	cc "github.com/ivanpirog/coloredcobra"
)

// Execute wires up colored help and runs the root command
func Execute() {
	root := NewRoot()

	cc.Init(&cc.Config{
		RootCmd:  root,
		Headings: cc.Bold + cc.Underline,
		Commands: cc.Bold + cc.HiMagenta,
		Example:  cc.Italic,
		ExecName: cc.Bold + cc.HiRed,
		Flags:    cc.Bold + cc.Yellow,
	})

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
