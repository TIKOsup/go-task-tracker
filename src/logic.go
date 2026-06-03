package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strconv"
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
	if desc == "" {
		log.Fatal("Task description cannot be empty")
	}

	fileData, err := os.ReadFile(DATA_FILE_PATH)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var tasks Tasks
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			log.Fatal("Error parsing JSON:", err)
		}
	}

	lastId, err := GetLastTaskId()
	if err != nil {
		log.Fatal("Error getting last task ID:", err)
	}

	newTask := Task{
		Id:          lastId + 1,
		Description: desc,
	}

	tasks.Tasks = append(tasks.Tasks, newTask)

	updatedJson, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	err = os.WriteFile(DATA_FILE_PATH, updatedJson, 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.Id)
}

func ListTasks(cmd *cobra.Command, args []string) {
	jsonFile, err := os.Open(DATA_FILE_PATH)
	if err != nil {
		log.Fatal("Error opening file:", err)
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
		return 0, nil
	}

	return tasks.Tasks[len(tasks.Tasks)-1].Id, nil
}

func DeleteTask(cmd *cobra.Command, args []string) {
	targetId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("Invalid task ID:", err)
	}

	fileData, err := os.ReadFile(DATA_FILE_PATH)
	if err != nil {
		log.Fatal("Error reading file:", err)
	}

	var tasks Tasks
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			log.Fatal("Error parsing JSON:", err)
		}
	}

	idx := slices.IndexFunc(tasks.Tasks, func(t Task) bool {
		return t.Id == targetId
	})

	if idx == -1 {
		log.Fatalf("Task with ID %d not found", targetId)
	}

	tasks.Tasks = slices.Delete(tasks.Tasks, idx, idx+1)

	updatedJson, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	err = os.WriteFile(DATA_FILE_PATH, updatedJson, 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	fmt.Printf("Task deleted successfully (ID: %d)\n", targetId)
}
