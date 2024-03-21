package main

import (
	"agen/task"
	"errors"
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
		listArgs := os.Args[2:]
		if len(listArgs) != 0 {
			if checkForHelpAndPrintUsage(listArgs, listUsage()) {
				os.Exit(0)
			}
		}
		tasks, err := task.LoadTasks()
		if err != nil {
			logAndExit(err.Error())
		}
		tasks, err = task.FilterTasks(tasks, listArgs)
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
// "done", the string slice can be empty and contains the uuids of part of it
// of the tasks to mark. Returns a non-nil error if the given tasks were marked.
func handleStatusMark(status string, args []string) error {
	stat, err := task.ParseStatus(status)
	if err != nil {
		return err
	}
	for _, uuid := range args {
		existsAndUnique, err := task.ExistsAndIsUnique(uuid)
		if err != nil {
			return err
		}
		if existsAndUnique {
			ts, err := task.LoadTask(uuid)
			if err != nil {
				return err
			}
			if err = ts.SetStatus(stat); err != nil {
				return err
			}
			if err = ts.SaveOnDisk(); err != nil {
				return err
			}
		} else {
			return errors.New("uuid prefix not unique")
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
		existsAndUnique, err := task.ExistsAndIsUnique(name)
		if err != nil {
			return err
		}
		if existsAndUnique {
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
		} else {
			return errors.New("uuid prefix not unique")
		}

	}
	return nil
}

// Removes the tasks denoted by the given uuids or part of it. If something
// wrong happens, returns an error. The args slice can be empty.
func handleRemove(args []string) error {
	for _, uuid := range args {
		if err := task.Remove(uuid); err != nil {
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
and t0 t1 ... denotes the optionnal tasks uuids (or part of it) to mark with
the given value`
}

// checks every element of args for equality with "-h", "--help" or "help" and
// if one is found that equals one of the strings, prints usage to the log
// and returns true
func checkForHelpAndPrintUsage(args []string, usage string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "help" || arg == "-help" {
			logger.Println(usage)
			return true
		}
	}
	return false
}

func listUsage() string {
	return `Usage of list:
  agen list [filter ...]
where filter is one of the following:
  status: todo, doing, done
  priority: low, medium, high

When several filters from the same category ("status" or "priority") are given,
they form a union filter, meaning that tasks that satisfy one of the given
filters could be listed (if not filtered out by the other category). If filters
from different categories are given, they form an intersection filter, meaning
that a task must have a status in the status filters and a priority in the
priority filters.

Examples:
  - to list all done tasks: agen list done
  - to list all done or todo tasks: agen list done todo
  - to list all todo tasks that have priority high: agen list todo high`
}
