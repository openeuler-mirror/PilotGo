package commands

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsPath string

func NewDocsCmd() *cobra.Command {
	var docsCmd = &cobra.Command{
		Use:     "docs",
		Short:   "generate pilotgo cmd docs",
		Example: `pilotgo docs`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return os.MkdirAll(docsPath, 0755)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return doc.GenMarkdownTree(rootCmd, docsPath)
		},
	}
	docsCmd.Flags().StringVarP(&docsPath, "path", "p", "./docs", "path to output docs")
	return docsCmd
}
