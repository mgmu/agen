package main

import (
	"agen/task"
	"flag"
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "agen:", log.LstdFlags)

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
	newTaskCmdTitle := newTaskCmd.String("title", "", `The task title.
Length must be strictly superior to 0 and strictly inferior to 256.`)
	newTaskCmdDesc := newTaskCmd.String("desc", "", `The task description.
Length must be strictly inferior to 65536.
This is optionnal and defaults to the empty string.`)
	newTaskCmdPeriod := newTaskCmd.Bool("periodic", false,
		`Indicates if the task is periodic.
This is optionnal and defaults to false.`)
	newTaskCmdPriority := newTaskCmd.Int("prio", int(task.Medium),
		`The task priority.
0 for Low, 1 for Medium and 2 for High.
This is optionnal and defaults to Medium.`)
	newTaskCmdStatus := newTaskCmd.Int("status", int(task.Todo),
		`The task status.
3 for Todo, 4 for Doing and 5 for Done.
This is optionnal and defaults to Todo.`)

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
		exists, err := task.Exists(ts.Title())
		if err != nil {
			logAndExit(err.Error())
		}
		if exists {
			fmt.Printf("Could not create task \"%s\": already exists\n",
				ts.Title())
			os.Exit(0)
		}
		if err = ts.SetDescription(*newTaskCmdDesc); err != nil {
			logAndExit(err.Error())
		}
		ts.SetPeriodicity(*newTaskCmdPeriod)
		if err = ts.SetPriority(byte(*newTaskCmdPriority)); err != nil {
			logAndExit(err.Error())
		}
		if err = ts.SetStatus(byte(*newTaskCmdStatus)); err != nil {
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
		for _, task := range tasks {
			fmt.Printf("> %s\n", task.Display())
		}
	default:
		logAndExit("unknown subcommand: " + os.Args[1])
	}
}
