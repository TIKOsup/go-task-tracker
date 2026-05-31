package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Task struct {
	Description string `json:"description"`
}

func AddTask(cmd *cobra.Command, args []string) {
	desc := strings.Join(args, " ")

	task := Task{
		Description: desc,
	}

	taskJson, _ := json.Marshal(task)

	fmt.Println(string(taskJson))

	err := os.WriteFile("./src/data.json", taskJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task added successfully (ID: TODO)")
}
