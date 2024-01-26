package main

import (
	"flag"
	"log"
	"os"
	"agen/task"
)

var logger = log.New(os.Stderr, "agen:", log.LstdFlags)

// Prints the given message on the logger and exits the program with exit status
// code 1
func logAndExit(msg string) {
	logger.Println(msg)
	os.Exit(1)
}

// Checks that the directory `$HOME/.agen/tasks` exists. Exits with status code
// 1 if it does not exist.
func checkTasksDirOrExit() {
	homePath := os.Getenv("HOME")
	if homePath == "" {
		logAndExit("$HOME not set")
	}
	f, err := os.Open(homePath + "/.agen/tasks")
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
		task, err := task.NewDefault(*newTaskCmdTitle)
		if err != nil {
			logAndExit(err.Error())
		}
		if err = task.SaveOnDisk(); err != nil {
			logAndExit(err.Error())
		}
	default:
		logAndExit("expected subcommand")
	}
}
