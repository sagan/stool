package versioncmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/sagan/stool/cmd"
	"github.com/sagan/stool/version"
)

var command = &cobra.Command{
	Use:   "version",
	Short: "Display program version.",
	Long:  `Display program version.`,
	Args:  cobra.MatchAll(cobra.ExactArgs(0), cobra.OnlyValidArgs),
	RunE:  versioncmd,
}

func init() {
	cmd.RootCmd.AddCommand(command)
}

func versioncmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("stool %s\n", version.Version)
	fmt.Printf("- build/date: %s\n", version.Date)
	fmt.Printf("- build/commit: %s\n", version.Commit)
	fmt.Printf("- os/type: %s\n", runtime.GOOS)
	fmt.Printf("- os/arch: %s\n", runtime.GOARCH)
	fmt.Printf("- go/version: %s\n", runtime.Version())
	return nil
}
