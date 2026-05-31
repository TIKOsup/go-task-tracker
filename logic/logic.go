package logic

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func AddTask(cmd *cobra.Command, args []string) {
	description := strings.Join(args, " ")
	fmt.Println("task added: " + description)
}
