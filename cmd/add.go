package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stikypiston/satchel/internal"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a file to the satchel",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		idx, err := internal.Load()
		if err != nil {
			return err
		}

		info, err := os.Stat(path)
		if err != nil {
			return err
		}

		id := idx.NextID
		dest := filepath.Join(
			internal.ObjectsDir(),
			fmt.Sprintf("%d", id),
		)

		if err := copyFile(path, dest); err != nil {
			return err
		}

		idx.Items[id] = internal.Item{
			Name:  filepath.Base(path),
			Size:  info.Size(),
			Added: info.ModTime(),
		}

		idx.NextID++

		if err := internal.Save(idx); err != nil {
			return err
		}

		fmt.Printf("Added '%s' as %d\n", path, id)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
