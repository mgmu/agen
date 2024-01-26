package task

import (
	"errors"
)

const (
	Low = iota
	Medium
	High
	Todo
	Doing
	Done
	TitleMinLength = 1
	TitleMaxLength = 255
)

var (
	ErrTitleTooShort   = errors.New("title too short")
	ErrTitleTooLong    = errors.New("title too long")
	ErrInvalidPriority = errors.New("priority must be Low, Medium or High")
	ErrInvalidStatus   = errors.New("status must be Todo, Doing or Done")
)

// A Task represents something to do before an arbitrary due date.
type Task struct {
	title      string // the title of the task (0 < length < 256)
	desc       string // the description of the task
	isPeriodic bool   // indicates if the task is periodic
	priority   byte   // the priority of the task: Low(0), Medium(1), High(2)
	status     byte   // the status of the task: Todo(3), Doing(4), Done(5)
}

func NewTask(title, desc string, isPeriodic bool, priority, status byte) (*Task,
	error) {
	if err := checkTitleValidity(title); err != nil {
		return nil, err
	}
	if priority > 2 {
		return nil, ErrInvalidPriority
	}
	if status < 3 || status > 5 {
		return nil, ErrInvalidStatus
	}
	new := Task{
		title:      title,
		desc:       desc,
		isPeriodic: isPeriodic,
		priority:   priority,
		status:     status,
	}
	return &new, nil
}

func NewDefault(title string) (*Task, error) {
	if err := checkTitleValidity(title); err != nil {
		return nil, err
	}
	new := Task{
		title:    title,
		priority: Medium,
		status:   Todo,
	}
	return &new, nil
}

// Returns the title of the task.
func (t *Task) Title() string {
	return t.title
}

// Returns the description of the task.
func (t *Task) Description() string {
	return t.desc
}

// Returns true if the task is periodic.
func (t *Task) IsPeriodic() bool {
	return t.isPeriodic
}

// Returns the priority of the task. The priority is one of these values: Low,
// Medium or High.
func (t *Task) Priority() int {
	return int(t.priority)
}

// Returns the status of the task. The status is one of these values: Todo,
// Doing or Done.
func (t *Task) Status() int {
	return int(t.status)
}

// Returns true if the given title is longer than the minimum title length
func isTitleLongerThanMinLength(title string) bool {
	return len(title) >= TitleMinLength
}

// Returns true if the given title is shorter than the maximum title length
func isTitleShorterThanMaxLength(title string) bool {
	return len(title) <= TitleMaxLength
}

// Returns nil if the title is valid, that is a string of
// length >= TitleMinLength and <= TitleMaxLength, an error otherwise
func checkTitleValidity(title string) error {
	if !isTitleLongerThanMinLength(title) {
		return ErrTitleTooShort
	}
	if !isTitleShorterThanMaxLength(title) {
		return ErrTitleTooLong
	}
	return nil
}
