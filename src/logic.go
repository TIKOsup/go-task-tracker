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
	"time"

	"github.com/spf13/cobra"
)

const DATA_FILE_PATH = "./src/data.json"
const STATUS_TODO = "todo"
const STATUS_IN_PROGRESS = "in-progress"
const STATUS_DONE = "done"

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

func AddTask(cmd *cobra.Command, args []string) {
	desc := strings.Join(args, " ")
	if desc == "" {
		log.Fatal("Task description cannot be empty")
	}

	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
	}

	lastId, err := GetLastTaskId()
	if err != nil {
		log.Fatal("Error getting last task ID:", err)
	}

	newTask := Task{
		Id:          lastId + 1,
		Description: desc,
		Status:      STATUS_TODO,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
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
	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
	}

	for _, task := range tasks.Tasks {
		fmt.Println(task)
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

	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
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

func GetData() (Tasks, error) {
	fileData, err := os.ReadFile(DATA_FILE_PATH)
	if err != nil {
		return Tasks{}, fmt.Errorf("error reading file: %w", err)
	}

	var tasks Tasks
	if len(fileData) > 0 {
		err = json.Unmarshal(fileData, &tasks)
		if err != nil {
			return Tasks{}, fmt.Errorf("error parsing JSON: %w", err)
		}
	}
	return tasks, nil
}

func UpdateTaskStatus(cmd *cobra.Command, args []string) {
	targetId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("Invalid task ID:", err)
	}

	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
	}

	idx := slices.IndexFunc(tasks.Tasks, func(t Task) bool {
		return t.Id == targetId
	})

	if idx == -1 {
		log.Fatalf("Task with ID %d not found", targetId)
	}

	var newStatus string
	switch cmd.CalledAs() {
	case "mark-in-progress":
		newStatus = STATUS_IN_PROGRESS
	case "mark-done":
		newStatus = STATUS_DONE
	default:
		log.Fatal("Unknown command:", cmd.CalledAs())
	}

	tasks.Tasks[idx].Status = newStatus

	updatedJson, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	err = os.WriteFile(DATA_FILE_PATH, updatedJson, 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	fmt.Printf("Task status updated successfully (ID: %d, New Status: %s)\n", targetId, newStatus)
}

func UpdateTaskDescription(cmd *cobra.Command, args []string) {
	targetId, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("Invalid task ID:", err)
	}

	newDescription := strings.Join(args[1:], " ")
	if newDescription == "" {
		log.Fatal("Description cannot be empty")
	}

	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
	}

	idx := slices.IndexFunc(tasks.Tasks, func(t Task) bool {
		return t.Id == targetId
	})

	if idx == -1 {
		log.Fatalf("Task with ID %d not found", targetId)
	}

	tasks.Tasks[idx].Description = newDescription

	updatedJson, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
	}

	err = os.WriteFile(DATA_FILE_PATH, updatedJson, 0644)
	if err != nil {
		log.Fatal("Error writing file:", err)
	}

	fmt.Printf("Task description updated successfully (ID: %d)\n", targetId)
}
