package main

import (
	"flag"
	"log"
	"os"
	"agen/task"
)

var logger = log.New(os.Stderr, "agen:", log.LstdFlags)

func main() {

	newTaskCmd := flag.NewFlagSet("newTask", flag.ExitOnError)
	newTaskCmdTitle := newTaskCmd.String("title", "",
		"the title of the task, of length > 0 and < 256")

	if len(os.Args) < 2 {
		logger.Println("expected subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "newTask":
		newTaskCmd.Parse(os.Args[2:])
		if *newTaskCmdTitle == "" {
			newTaskCmd.Usage()
			os.Exit(1)
		}
		logger.Println("length of title", len(*newTaskCmdTitle))
		_, err := task.NewDefault(*newTaskCmdTitle)
		if err != nil {
			logger.Println("error")
			logger.Fatal(err.Error())
		}
	default:
		logger.Println("expected subcommand")
		os.Exit(1)
	}
}
