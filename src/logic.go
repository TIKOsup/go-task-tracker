package logic

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const dataFilePath = "data.json"
const statusTodo = "todo"
const statusInProgress = "in-progress"
const statusDone = "done"

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
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

	var lastId int
	if len(tasks.Tasks) == 0 {
		lastId = 0
	} else {
		lastId = tasks.Tasks[len(tasks.Tasks)-1].Id
	}

	newTask := Task{
		Id:          lastId + 1,
		Description: desc,
		Status:      statusTodo,
		CreatedAt:   time.Now(),
	}

	tasks.Tasks = append(tasks.Tasks, newTask)

	err = SaveData(tasks)
	if err != nil {
		log.Fatal("Error saving data:", err)
	}

	fmt.Printf("Task added successfully (ID: %d)\n", newTask.Id)
}

func ListTasks(cmd *cobra.Command, args []string) {
	tasks, err := GetData()
	if err != nil {
		log.Fatal("Error getting data:", err)
	}

	if len(args) == 1 {
		filter := args[0]

		if filter != statusTodo && filter != statusInProgress && filter != statusDone {
			log.Fatalf("Invalid status filter: %s. Valid options are: %s, %s, %s", filter, statusTodo, statusInProgress, statusDone)
		}

		filteredTasks := make([]Task, 0, len(tasks.Tasks))
		for _, t := range tasks.Tasks {
			if t.Status == filter {
				filteredTasks = append(filteredTasks, t)
			}
		}
		tasks.Tasks = filteredTasks
	}

	for _, task := range tasks.Tasks {
		fmt.Printf("[%d] %s (%s)\n", task.Id, task.Description, task.Status)
	}
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

	err = SaveData(tasks)
	if err != nil {
		log.Fatal("Error saving data:", err)
	}

	fmt.Printf("Task deleted successfully (ID: %d)\n", targetId)
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
		newStatus = statusInProgress
	case "mark-done":
		newStatus = statusDone
	default:
		log.Fatal("Unknown command:", cmd.CalledAs())
	}

	tasks.Tasks[idx].Status = newStatus
	tasks.Tasks[idx].UpdatedAt = time.Now()

	err = SaveData(tasks)
	if err != nil {
		log.Fatal("Error saving data:", err)
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
	tasks.Tasks[idx].UpdatedAt = time.Now()

	err = SaveData(tasks)
	if err != nil {
		log.Fatal("Error saving data:", err)
	}

	fmt.Printf("Task description updated successfully (ID: %d)\n", targetId)
}

func GetData() (Tasks, error) {
	fileData, err := os.ReadFile(dataFilePath)

	if os.IsNotExist(err) {
		return Tasks{
			Tasks: []Task{},
		}, nil
	}

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

func SaveData(tasks Tasks) error {
	updatedJson, err := json.MarshalIndent(tasks, "", " ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %w", err)
	}

	err = os.WriteFile(dataFilePath, updatedJson, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %w", err)
	}
	return nil
}
