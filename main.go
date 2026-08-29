// Copyright (c) Alex Ellis 2023. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"log"
	"os"

	"github.com/alexellis/fstail/pkg/fstail"
	"github.com/spf13/cobra"
)

func main() {
	var prefixStyle string
	var disablePrefix bool

	var rootCmd = &cobra.Command{
		Use:   "fstail",
		Short: "Tail files with optional string matching",
		Long: `Tail files in a directory with optional string matching and prefix styles.

fstail uses fsnotify to detect both new files, and existing files that are changed.
It then starts concatenating their contents to the terminal.`,
		Example: `  # Work in current directory
  fstail

  # Work in /var/log/actuated
  fstail /var/log/actuated

  # Work in /var/log/actuated and match strings with "server error"
  fstail /var/log/actuated "server error"

  # Disable prefix printing
  fstail --prefix none /var/log/actuated

  # Print a prefix with the calculated Pod/container name
  fstail --prefix k8s /var/log/containers "queue-worker"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var opts fstail.RunOptions

			if len(args) == 1 {
				opts.WorkDir = args[0]
			} else if len(args) == 2 {
				opts.WorkDir = args[0]
				opts.Match = args[1]
			} else {
				cwd, err := os.Getwd()
				if err != nil {
					return err
				}
				opts.WorkDir = cwd
			}

			opts.DisableLogPrefix = disablePrefix
			if disablePrefix || prefixStyle == "none" {
				opts.PrefixStyle = fstail.PrefixStyleNone
			} else if prefixStyle == "k8s" {
				opts.PrefixStyle = fstail.PrefixStyleK8s
			} else {
				opts.PrefixStyle = fstail.PrefixStyleFilename
			}

			return fstail.Run(opts)
		},
	}

	rootCmd.Flags().StringVar(&prefixStyle, "prefix", "filename", "Prefix style: filename, k8s, or none")
	rootCmd.Flags().BoolVar(&disablePrefix, "no-prefix", false, "Disable prefix printing (equivalent to --prefix none)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
