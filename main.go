package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"task-cli/handler"
)

func main() {

	fmt.Println("---------TASK TRACKER CLI---------")
	for {
		fmt.Print(">task-cli ")
		var input = bufio.NewScanner(os.Stdin)
		input.Scan()
		fields := strings.Fields(input.Text())

		cmd := fields[0]
		// fmt.Println(cmd, fields)
		switch cmd {
		case "add":

			if len(fields) < 2 {
				log.Fatalf("insufficient arguments, expected 1, received zero")
			}
			handler.HandleAdd(fields[1:])
		case "list":

			handler.HandleList(fields[1:])
		case "mark-in-progress", "mark-done":
			handler.HandleMark(fields)
		case "exit":
			os.Exit(0)
		default:
			log.Fatalf("invalid command. %v\n", cmd)
			return
		}
	}

}
