package list

import (
	"github.com/spf13/cobra"
)

var CMDList = &cobra.Command{
	Use:   "list",
	Short: "List information from the API",
}
