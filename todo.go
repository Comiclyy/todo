package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"
	"time"
)

const tasksFile = "/home/max/.config/todo/tasks.json"

var tasks []string

func main() {
	loadTasks()

	if len(os.Args) > 1 {
		command := os.Args[1]

		switch command {
		case "list":
			list()
		case "add":
			add()
		case "remove":
			remove()
		case "complete":
			complete()
		default:
			fmt.Println("Invalid command. Available commands: list, add, remove, complete")
		}
	} else {
		fmt.Println("What would you like to do:\nList tasks(L)\nAdd task(A)\nRemove task(R)\nComplete task(C)")
		var input string
		fmt.Scanln(&input)

		switch strings.ToLower(input) {
		case "l":
			list()
		case "a":
			add()
		case "r":
			remove()
		case "c":
			complete()
		default:
			fmt.Println("Invalid input. Please try again.")
		}
	}

	saveTasks()
}

func add() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the task:")

	task, _ := reader.ReadString('\n')
	task = strings.TrimSpace(task)

	if task != "" {
		tasks = append(tasks, task)
	}

	list()
}

func remove() {
	if len(tasks) == 0 {
		fmt.Println("No tasks to remove.")
		return
	}

	list()

	fmt.Println("Enter the task number to remove:")

	var input int
	_, err := fmt.Scanln(&input)
	if err != nil || input < 1 || input > len(tasks) {
		fmt.Println("Invalid task number. Please try again.")
		return
	}

	removedTask := tasks[input-1]
	tasks = append(tasks[:input-1], tasks[input:]...)

	fmt.Printf("Removed task: %s\n", removedTask)
}

func complete() {
	if len(tasks) == 0 {
		fmt.Println("No tasks to complete.")
		return
	}

	list()

	fmt.Println("Enter the task number to complete:")

	var input int
	_, err := fmt.Scanln(&input)
	if err != nil || input < 1 || input > len(tasks) {
		fmt.Println("Invalid task number. Please try again.")
		return
	}

	task := tasks[input-1]
	if strings.Contains(task, "✓") {
		fmt.Println("Task is already complete.")
		return
	}

	tasks[input-1] = task + " ✓"

	fmt.Printf("Completed task: %s\n", task)
}

func list() {
	currentUser, err := user.Current()
	if err != nil {
		fmt.Println("Failed to get current user:", err)
		return
	}

	currentTime := time.Now().Format("15:04:05")

	yellow := "\033[33m"
	purple := "\033[35m"
	green := "\033[32m"
	white := "\033[37m"
	red := "\033[31m"
	bold := "\033[1m"
	reset := "\033[0m"

	lineLength := 18
	line := strings.Repeat("⎯", lineLength)

	greeting := fmt.Sprintf("%s Hello %s, it's %s %s", yellow+line+reset, currentUser.Username, currentTime, yellow+line+reset)

	fmt.Printf("%-42s %s\n", greeting, "")

	fmt.Printf("%-5s %-25s %s\n", yellow+"ID", "   Task", "   Status")

	for i, task := range tasks {
		taskNumber := fmt.Sprintf("%s%-5d", purple, i+1)
		taskText := fmt.Sprintf("%-25s", task)
		isComplete := strings.Contains(task, "✓")
		status := ""

		if isComplete {
			taskText = green + bold + taskText + reset
			status = green + bold + "✓" + reset
		} else {
			taskText = white + taskText + reset
			status = red + bold + "⃝" + reset
		}

		fmt.Printf("%s %s %s\n", taskNumber, taskText, status)
	}
}

func loadTasks() {
	if _, err := os.Stat(tasksFile); err == nil {
		content, err := ioutil.ReadFile(tasksFile)
		if err == nil {
			err = json.Unmarshal(content, &tasks)
			if err != nil {
				fmt.Println("Failed to load tasks:", err)
			}
		}
	}
}

func saveTasks() {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Failed to encode tasks:", err)
		return
	}

	err = ioutil.WriteFile(tasksFile, data, 0644)
	if err != nil {
		fmt.Println("Failed to save tasks:", err)
	}
}
