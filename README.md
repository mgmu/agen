# agen
A simple todo list manager CLI application.

# Install
First of all, run `./install.sh`. It will create a hidden directory in your
`$HOME` named `.agen` where the tasks will be stored, then run `go build` in
order to generate the `agen` binary.

# Uninstall
Just remove the directory `$HOME/.agen/` and the binary.

# Usage
Every task has a mandatory title that must be provided at task creation. The
rest can be left blank and default values will be provided. Also, a unique
indentifier is given to every created task. This identifier is used to edit the
task, by mark or remove. To add a new task with title "Prep dinner" and default
parameters, run:  
`
agen newTask -title "Prep dinner"
`
  
This will create a new task, not periodic, without a description, of status
"To do", of medium priority and a unique identifier.  
  
To list all the tasks, run:  
`
agen list
`
  
To change the status or the priority of tasks, run `agen mark` followed by the
value you want to set and a list that can be empty of task identifiers. This
list can contain full identifiers, but as their are rather long to type, you can
type only the beginning of the identifiers, the programs looks for the task that
has the corresponding identifier. If multiple tasks exists with the same
identifier prefix, the modification is not applied, to none of the corresponding
tasks in the list. For example, if you created a task and it has an identifier
that starts with 3a and its the only task that starts with that prefix, you can
change its priority to high with:  
`
agen mark high 3a
`
  
If you prepared the dinner, run:  
`
agen mark done 3a
`

If you don't want to hear about the fact that you have to prepare dinner
anymore, run:  
`
agen remove 3a
`
