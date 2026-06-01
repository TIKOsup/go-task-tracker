package logic

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
}

func AddTask(cmd *cobra.Command, args []string) {
	desc := strings.Join(args, " ")

	task := Task{
		Id:          0,
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

func ListTasks(cmd *cobra.Command, args []string) {
	jsonFile, err := os.Open("./src/data.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var tasks Tasks
	err = json.Unmarshal(byteValue, &tasks)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	for _, task := range tasks.Tasks {
		fmt.Println("ID:", task.Id, "Description:", task.Description)
	}
}
