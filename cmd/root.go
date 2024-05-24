/*
Copyright © 2023 Alec DuBois <alec@petrichormud.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ptcr",
	Short: "The Petrichor App",
	Long:  `The Petrichor App`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
