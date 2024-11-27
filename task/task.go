package task

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type TaskStatus string

type Tasks struct {
	Tasks []Task `json:"tasks"`
}

const (
	TODO       TaskStatus = "todo"
	INPROGRESS TaskStatus = "in-progress"
	DONE       TaskStatus = "done"
)

var taskId int = 0
var filePath string = "tasks.json"
var tasks = Tasks{
	Tasks: make([]Task, 0),
}
var ValidTaskStatus = map[TaskStatus]bool{
	TODO:       true,
	INPROGRESS: true,
	DONE:       true,
}

func (t Task) String() string {
	return fmt.Sprintf("\n {\n  Id: %d\n  Description: %s\n  Status: %s\n  CreatedAt: %v\n  UpdatedAt: %v\n }\n", t.Id, t.Description, t.Status, t.CreatedAt, t.UpdatedAt)
}

func NewTask(desc string) Task {
	taskId++
	t := Task{
		Id:          taskId,
		Description: desc,
		Status:      TODO,
		CreatedAt:   time.Now(),
	}
	tasks.Tasks = append(tasks.Tasks, t)

	return t
}

func init() {
	// var tasksTmp Tasks
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		tasks.Tasks = make([]Task, 0)
		return
	}
	buf, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("Error reading file %s: %v", filePath, err)
	}

	err = json.Unmarshal(buf, &tasks)

	if err != nil {
		log.Fatalf("error while decoding json: %v", err)
	}

	// tasks.Tasks = tasksTmp.Tasks
	taskId = len(tasks.Tasks)
}

func Save() error {

	log.Printf("Tasks: %v", tasks.Tasks)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return fmt.Errorf("error while opening file: %v", err)
	}

	defer f.Close()

	buf, err := json.MarshalIndent(tasks, "", " ")

	if err != nil {
		return fmt.Errorf("error while encoding json: %v", err)
	}

	_, err = f.Write(buf)
	if err != nil {
		return fmt.Errorf("error while writing to file: %v", err)
	}

	return nil
}

func MarkTask(taskId int, status string) (int, error) {
	for i := range tasks.Tasks {
		if tasks.Tasks[i].Id == taskId {
			tasks.Tasks[i].Status = markStatus(status)
			tasks.Tasks[i].UpdatedAt = time.Now()
			log.Printf("Updated Tasks: %v", tasks.Tasks)
			err := Save()
			if err != nil {
				return 0, err
			}
			return taskId, nil
		}
	}

	return 0, fmt.Errorf("task with id %d not found", taskId)
}

func markStatus(status string) TaskStatus {
	if status == "mark-in-progress" {
		return INPROGRESS
	} else if status == "mark-done" {
		return DONE
	}
	fmt.Println("Invalid status. Use 'mark-in-progress' or 'mark-done'")
	return TODO

}

func FetchAll() (*[]Task, error) {

	// var tasksTmp Tasks

	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return nil, fmt.Errorf("list does not exist: %v", err)
	}

	buf, err := os.ReadFile(filePath)

	if err != nil {
		return nil, fmt.Errorf("error while reading file: %v", err)
	}

	err = json.Unmarshal(buf, &tasks)

	if err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}

	return &tasks.Tasks, nil
}

func FetchByStatus(status TaskStatus) (*[]Task, error) {

	_, err := FetchAll()

	if err != nil {
		return nil, err
	}
	filteredTasks := filter(func(task Task) bool { return task.Status == status })

	return &filteredTasks, nil
}

func FetchById(id int) (*Task, error) {

	// _, err := FetchAll()

	// if err != nil {
	// 	return &Task{}, err
	// }
	// log.Printf("Address of task: %p", &tasks)
	filteredTasks := filter(func(task Task) bool { return task.Id == id })

	return &filteredTasks[0], nil
}

func filter(f func(task Task) bool) []Task {

	result := []Task{}
	for _, task := range tasks.Tasks {
		if f(task) {
			result = append(result, task)
		}
	}
	return result
}

func UpdateTask(id int, description string) error {
	for i := range tasks.Tasks {
		if tasks.Tasks[i].Id == id {
			tasks.Tasks[i].Description = description
			tasks.Tasks[i].UpdatedAt = time.Now()
			err := Save()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return fmt.Errorf("task with id %d not found", id)
}
