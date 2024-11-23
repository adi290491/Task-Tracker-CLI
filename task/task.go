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

var taskId int = 1
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
	t := Task{
		Id:          taskId,
		Description: desc,
		Status:      TODO,
		CreatedAt:   time.Now(),
	}
	tasks.Tasks = append(tasks.Tasks, t)
	taskId++
	return t
}

// func (t *Task) UnmarshalJSON(p []byte) error {

// 	log.Println("Calling Unmarshal JSON as intermediate")
// 	var tmp []interface{}

// 	if err := json.Unmarshal(p, &tmp); err != nil {
// 		return err
// 	}

// 	t.Id = tmp[0].(int)
// 	t.Description = tmp[1].(string)
// 	t.Status = TaskStatus(tmp[2].(string))
// 	t.CreatedAt = tmp[3].(time.Time)
// 	t.UpdatedAt = tmp[4].(time.Time)

// 	return nil
// }

func (t *Task) Save() (int, error) {

	log.Printf("Tasks: %v", tasks.Tasks)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

	if err != nil {
		return 0, fmt.Errorf("error while opening file: %v", err)
	}

	defer f.Close()

	buf, err := json.MarshalIndent(tasks, "", " ")

	if err != nil {
		return 0, fmt.Errorf("error while encoding json: %v", err)
	}

	_, err = f.Write(buf)
	if err != nil {
		return 0, fmt.Errorf("error while writing to file: %v", err)
	}

	return t.Id, nil
}

func MarkTask(taskId int, status string) (int, error) {
	task, err := FetchById(taskId)

	if err != nil {
		return 0, fmt.Errorf("task with id %d not found", taskId)
	}

	task.Status = markStatus(status)
	task.UpdatedAt = time.Now()
	log.Printf("Found Task: %v", *task)

	return task.Save()
}

// func MarkDone(taskId int) (int, error) {

// 	for _, task := range tasks.Tasks {
// 		if task.Id == taskId {
// 			task.Status = DONE
// 			task.UpdatedAt = time.Now()
// 			return task.Save()
// 		}
// 	}
// 	return 0, fmt.Errorf("task with id %d not found", taskId)
// }

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

	var tasks Tasks

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

	tasks, err := FetchAll()

	if err != nil {
		return nil, err
	}
	filteredTasks := filter(*tasks, func(task Task) bool { return task.Status == status })

	return &filteredTasks, nil
}

func FetchById(id int) (*Task, error) {

	tasks, err := FetchAll()

	if err != nil {
		return &Task{}, err
	}
	filteredTasks := filter(*tasks, func(task Task) bool { return task.Id == id })

	return &filteredTasks[0], nil
}

func filter(tasks []Task, f func(task Task) bool) []Task {
	result := []Task{}
	for _, task := range tasks {
		if f(task) {
			result = append(result, task)
		}
	}
	return result
}
