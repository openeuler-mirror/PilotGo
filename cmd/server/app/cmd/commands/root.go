/*
 * Copyright (c) KylinSoft  Co., Ltd. 2024.All rights reserved.
 * PilotGo licensed under the Mulan Permissive Software License, Version 2.
 * See LICENSE file for more details.
 * Author: linjieren <linjieren@kylinos.cn>
 * Date: Thu Jul 25 16:18:53 2024 +0800
 */
package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pilotgo",
	Short: "PilotgGo CLI",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
	DisableAutoGenTag: true,
	SilenceUsage:      true,
}

func Execute() {
	if len(os.Args) == 1 {
		rootCmd.SetArgs([]string{cliName})
	}
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewServerCommand())
	rootCmd.AddCommand(NewTempleteCommand())
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewDocsCmd())
}
