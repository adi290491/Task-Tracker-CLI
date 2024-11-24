package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"task-cli/task"
)

func HandleAdd(args []string) {

	desc := strings.Join(args[0:], " ")
	desc = strings.ReplaceAll(desc, "\"", "")

	var task = task.NewTask(desc)
	id, err := task.Save()
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("Task added successfully (ID: %d)\n", id)

}

func HandleList(status []string) {
	file, err := os.Open("tasks.json")
	if err != nil {
		log.Fatalf("No tasks created for status")
		return
	}
	defer file.Close()

	if len(status) == 0 {

		_, err := task.FetchAll()

		if err != nil {
			log.Fatalf("Error fetching tasks: %v", err)
		}

		// fmt.Println(tasks)

	} else {

		if !task.ValidTaskStatus[task.TaskStatus(status[0])] {
			log.Fatalf("Undefined task status: %v", status)
		}

		task, err := task.FetchByStatus(task.TaskStatus(status[0]))

		if err != nil {
			log.Fatalf("Error fetching tasks: %v", err)
		}

		fmt.Println(*task)

	}

}

func HandleMark(cmd []string) {
	if len(cmd) < 2 {
		log.Fatalf("Insufficient arguments, expected Id, received zero")
	}
	markStatus := cmd[0]
	id, err := strconv.Atoi(cmd[1])

	if err != nil {
		log.Fatalf("Invalid ID: %v", err)
	}

	_, err = task.MarkTask(id, markStatus)
	if err != nil {
		log.Fatalf("Error marking task: %v", err)
	}
}
