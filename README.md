# agen
A simple todo list app.

# Install
First of all, run `./install.sh`. It will create a hidden directory in your  
`$HOME` named `.agen` where the tasks will be stored, then run `go build` in  
order to generate the `agen` binary.

# Uninstall
Just remove the directory `$HOME/.agen/` and the binary from where you put it.

# Usage
Every task has a mandatory title.  
To add a new task with title "Prep dinner" and default parameters, run:  
`
agen newTask -title "Prep dinner"
`
  
This will create a new task, not periodic, without a description, of status  
"To do" and medium priority.  
  
To list all the tasks, run:  
`
agen list
`
  
To change status or priority of tasks, run `agen mark` followed by the value  
you want to set. For example, if you really have to prepare the dinner later,  
run:  
`
agen mark high "Prep dinner"
`
  
If you prepared the dinner, run:  
`
agen mark done "Prep dinner"
`

If you don't want to hear about the fact that you have to prepare dinner  
anymore, run:  
`
agen remove "Prep dinner"
`
