package main

import (
	"agen/task"
	"flag"
	"fmt"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "agen:", log.LstdFlags)

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
	newTaskCmdPriority := newTaskCmd.String("prio", "medium",
		`The task priority.
"low" for Low, "medium" for Medium and "high" for High.
This is optionnal and defaults to Medium.`)
	newTaskCmdStatus := newTaskCmd.String("status", "todo",
		`The task status.
"todo" for Todo, "doing" for Doing and "done" for Done.
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
		prio, err := task.ParsePriority(*newTaskCmdPriority)
		if err != nil {
			logAndExit(err.Error())
		}
		if err = ts.SetPriority(prio); err != nil {
			logAndExit(err.Error())
		}
		status, err := task.ParseStatus(*newTaskCmdStatus)
		if err != nil {
			logAndExit(err.Error())
		}
		if err = ts.SetStatus(status); err != nil {
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
	case "mark":
		if len(os.Args) < 3 {
			logAndExit("no specific mark given")
		}
		switch os.Args[2] {
		case "todo", "doing", "done":
			if len(os.Args[2:]) < 2 {
				os.Exit(0)
			}
			if err := handleStatusMark(os.Args[2], os.Args[3:]); err != nil {
				logAndExit(err.Error())
			}
		case "low", "medium", "high":
			if len(os.Args[2:]) < 2 {
				os.Exit(0)
			}
			if err := handlePriorityMark(os.Args[2], os.Args[3:]); err != nil {
				logAndExit(err.Error())
			}
		default:
			logger.Println("unkown mark: " + os.Args[2])
			fmt.Println(markSubcommandUsage())
			os.Exit(1)
		}
	case "remove":
		if len(os.Args) < 3 {
			os.Exit(0)
		}
		if err := handleRemove(os.Args[2:]); err != nil {
			logAndExit(err.Error())
		}
	default:
		logAndExit("unknown subcommand: " + os.Args[1])
	}
}

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

// Handle for status marking, the given status must be either "todo", "doing" or
// "done", the string slice can be empty and contains the name of the tasks to
// mark. Returns a non-nil error if the given task names were marked.
func handleStatusMark(status string, args []string) error {
	stat, err := task.ParseStatus(status)
	if err != nil {
		return err
	}
	for _, name := range args {
		exists, err := task.Exists(name)
		if err != nil {
			return err
		}
		if exists {
			ts, err := task.LoadTask(name)
			if err != nil {
				return err
			}
			if err = ts.SetStatus(stat); err != nil {
				return err
			}
			if err = ts.SaveOnDisk(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Handle for priority marking, the given priority must be either "low",
// "medium" or "high", the string slice can be empty and contains the name of
// the tasks to mark. Returns a non-nil error if the given task names were
// marked.
func handlePriorityMark(priority string, args []string) error {
	prio, err := task.ParsePriority(priority)
	if err != nil {
		return err
	}
	for _, name := range args {
		exists, err := task.Exists(name)
		if err != nil {
			return err
		}
		if exists {
			ts, err := task.LoadTask(name)
			if err != nil {
				return err
			}
			if err = ts.SetPriority(prio); err != nil {
				return err
			}
			if err = ts.SaveOnDisk(); err != nil {
				return err
			}			
		}
	}
	return nil
}

// Removes the tasks denotes by the given names. If something wrong happens,
// returns an error. The args slice can be empty
func handleRemove(args []string) error {
	for _, name := range args {
		if err := task.Remove(name); err != nil { // todo
			return err
		}
	}
	return nil
}

func markSubcommandUsage() string {
	return `Usage of mark:
  agen mark arg [t0 t1 ...]
where arg is one of the following:
  low:    sets the priority of the given tasks to Low
  medium: sets the priority of the given tasks to Medium
  high:   sets the priority of the given tasks to High

  todo:   sets the status of the given tasks to Todo
  doing:  sets the status of the given tasks to Doing
  done:   sets the status of the given tasks to Done
and t0 t1 ... denotes the optionnal list of tasks to mark with the given value`
}
