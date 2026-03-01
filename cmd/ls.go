package cmd

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/stikypiston/satchel/internal"
)

// lsCmd represents the ls command
var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List the contents of the satchel",
	RunE: func(cmd *cobra.Command, args []string) error {
		idx, err := internal.Load()
		if err != nil {
			return err
		}

		if len(idx.Items) == 0 {
			fmt.Println("Satchel is empty")
			return nil
		}

		// Sort ID's
		var ids []int
		for id := range idx.Items {
			ids = append(ids, int(id))
		}
		sort.Ints(ids)

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"ID", "Name", "Size", "Added"})

		for _, id := range ids {
			item := idx.Items[int64(id)]

			table.Append([]string{
				strconv.Itoa(id),
				item.Name,
				humanSize(item.Size),
				timeAgo(item.Added),
			})
		}

		table.Render()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func humanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(size)/float64(div),
		"KMGTPE"[exp],
	)
}

func timeAgo(t time.Time) string {
	d := time.Since(t)

	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
