package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"task-cli/task"
)

func HandleAdd(fields []string) {
	if len(fields) < 2 {
		log.Println("insufficient arguments, expected 1, received zero")
		return
	}

	desc := strings.Join(fields[1:], " ")
	desc = strings.ReplaceAll(desc, "\"", "")

	var t = task.NewTask(desc)
	err := task.Save()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Task added successfully (ID: %d)\n", t.Id)

}

func HandleList(status []string) {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Printf("No tasks created for status")
		return
	}
	defer file.Close()

	if len(status) == 0 {

		tasks, err := task.FetchAll()

		if err != nil {
			log.Printf("Error fetching tasks: %v", err)
		}

		fmt.Println(*tasks)

	} else {

		if !task.ValidTaskStatus[task.TaskStatus(status[0])] {
			log.Printf("Undefined task status: %v", status)
		}

		task, err := task.FetchByStatus(task.TaskStatus(status[0]))

		if err != nil {
			log.Printf("Error fetching tasks: %v", err)
		}

		fmt.Println(*task)

	}

}

func HandleMark(cmd []string) {
	if len(cmd) < 2 {
		log.Printf("Insufficient arguments, expected Id, received zero")
	}
	markStatus := cmd[0]
	id, err := strconv.Atoi(cmd[1])

	if err != nil {
		log.Printf("Invalid ID: %v", err)
	}

	_, err = task.MarkTask(id, markStatus)
	if err != nil {
		log.Printf("Error marking task: %v", err)
	}
}

func HandleUpdate(fields []string) {
	if len(fields) < 2 {
		log.Printf("Insufficient arguments,\nUsage: update <id> <description>\n")
		return
	}

	id, _ := strconv.Atoi(fields[0])
	desc := strings.Join(fields[1:], " ")
	desc = strings.ReplaceAll(desc, "\"", "")

	err := task.UpdateTask(id, desc)
	if err != nil {
		log.Printf("Error updating task: %v\n", err)
	} else {
		fmt.Printf("Task updated successfully (ID: %d)\n", id)
	}
}

func HandleDelete(fields []string) {
	if len(fields) < 2 {
		log.Printf("Insufficient arguments,\nUsage: update <id>\n")
		return
	}

	id, _ := strconv.Atoi(fields[1])

	err := task.DeleteTask(id)

	if err == nil {
		fmt.Printf("Task deleted successfully (ID: %d)\n", id)
	} else {
		log.Printf("Error deleting task: %v\n", err)
	}

}
