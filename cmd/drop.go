package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/stikypiston/satchel/internal"
)

var dropYes bool

var dropCmd = &cobra.Command{
	Use:   "drop",
	Short: "Clear the entire satchel",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !dropYes {
			confirm := false

			ask := huh.NewConfirm().
				Title("Drop satchel?").
				Affirmative("Yes").
				Negative("No").
				Value(&confirm)
			ask.Run()

			if !confirm {
				fmt.Println("Aborted.")
				return nil
			}
		}

		dir := internal.BaseDir()
		if err := os.RemoveAll(dir); err != nil {
			return err
		}

		fmt.Println("Satchel dropped.")
		return nil
	},
}

func init() {
	dropCmd.Flags().BoolVarP(&dropYes, "yes", "y", false, "skip confirmation prompt")
	rootCmd.AddCommand(dropCmd)
}
