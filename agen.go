package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"agen/task"
)

var	logger = log.New(os.Stderr, "agen:", log.LstdFlags)

// Prints the given message on the logger and exits the program with exit status
// code 1
func logAndExit(msg string) {
	logger.Println(msg)
	os.Exit(1)
}

// Checks that the directory `$HOME/.agen/tasks` exists. Exits with status code
// 1 if it does not exist. Also sets the tasks save path
func checkTasksDirOrExit() {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		logAndExit("$HOME not set")
	}
	task.TasksPath = homePath + "/.agen/tasks"
	f, err := os.Open(task.TasksPath)
	if err != nil {
		logAndExit(err.Error())
	}
	fi, err := f.Stat()
	if err != nil {
		logAndExit(err.Error())
	}
	if !fi.IsDir() {
		logAndExit("tasks directory missing")
	}
}

func main() {
	checkTasksDirOrExit()

	newTaskCmd := flag.NewFlagSet("newTask", flag.ExitOnError)
	newTaskCmdTitle := newTaskCmd.String("title", "",
		"the title of the task, of length > 0 and < 256")

	_ = flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		logAndExit("expected subcommand")
	}

	switch os.Args[1] {
	case "newTask":
		newTaskCmd.Parse(os.Args[2:])
		if *newTaskCmdTitle == "" {
			newTaskCmd.Usage()
			os.Exit(1)
		}
		ts, err := task.NewDefault(*newTaskCmdTitle)
		if err != nil {
			logAndExit(err.Error())
		}
		if err = ts.SaveOnDisk(); err != nil {
			logAndExit(err.Error())
		}
	case "list":
		tasks, err := task.LoadTasks()
		if err != nil {
			logAndExit(err.Error())
		}
		if len(tasks) == 0 {
			fmt.Println("Nothing to show")
		}
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Title())
		}
	default:
		logAndExit("expected subcommand")
	}
}
