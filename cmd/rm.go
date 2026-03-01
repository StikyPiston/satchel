package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stikypiston/satchel/internal"
)

var rmYes bool

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm <id>",
	Short: "Remove a file from the satchel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return errors.New("invalid ID")
		}

		idx, err := internal.Load()
		if err != nil {
			return err
		}

		item, ok := idx.Items[id]
		if !ok {
			return errors.New("ID not found")
		}

		if !rmYes {
			fmt.Printf("Are you sure you want to delete %s (ID %d)? [y/N] ", item.Name, id)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("Aborted.")
				return nil
			}
		}

		// Remove object file
		objPath := filepath.Join(internal.ObjectsDir(), strconv.FormatInt(id, 10))
		if err := os.Remove(objPath); err != nil && !os.IsNotExist(err) {
			return err
		}

		// Remove from index
		delete(idx.Items, id)

		if err := internal.Save(idx); err != nil {
			return err
		}

		fmt.Printf("Removed %s (ID %d)\n", item.Name, id)
		return nil
	},
}

func init() {
	rmCmd.Flags().BoolVarP(&rmYes, "yes", "y", false, "Skip confirmation prompt")
	rootCmd.AddCommand(rmCmd)
}
