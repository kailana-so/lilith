package markdown

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/mattn/go-isatty"
)

func CLIFormatter(s string) {
	// If not a TTY (e.g., piped), print raw Markdown (better for | tee file.md).
	if !(isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		fmt.Println(s)
		return
	}

	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),         // picks light/dark based on $COLORTERM etc.
		glamour.WithWordWrap(100),       // adjust to your preferred width
		glamour.WithEnvironmentConfig(), // honors GLAMOUR_* env vars
	)
	if err != nil {
		fmt.Println(s) // graceful fallback
		return
	}

	out, err := r.Render(s)
	if err != nil {
		fmt.Println(s)
		return
	}
	fmt.Print(out)
}
