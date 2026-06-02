package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const DATA_FILE_PATH = "./src/data.json"

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

	err := os.WriteFile(DATA_FILE_PATH, taskJson, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Task added successfully (ID: TODO)")
}

func ListTasks(cmd *cobra.Command, args []string) {
	jsonFile, err := os.Open(DATA_FILE_PATH)
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

func GetLastTaskId() (int, error) {
	jsonFile, err := os.Open(DATA_FILE_PATH)
	if err != nil {
		return 0, errors.New("error opening file: " + err.Error())
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return 0, errors.New("error reading file: " + err.Error())
	}

	var tasks Tasks
	err = json.Unmarshal(byteValue, &tasks)
	if err != nil {
		return 0, errors.New("error parsing JSON: " + err.Error())
	}

	if len(tasks.Tasks) == 0 {
		return 0, errors.New("no tasks found")
	}

	return tasks.Tasks[len(tasks.Tasks)-1].Id, nil
}
