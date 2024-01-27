package task

import (
	"errors"
	"path/filepath"
	"os"
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
	DescMaxLength  = 65535
)

var (
	TasksPath = ""

	ErrTitleTooShort   = errors.New("title too short (min 1)")
	ErrTitleTooLong    = errors.New("title too long (max 255)")
	ErrInvalidPriority = errors.New("priority must be Low, Medium or High")
	ErrInvalidStatus   = errors.New("status must be Todo, Doing or Done")
	ErrDescTooLong     = errors.New("description too long (max 65535)")
)

// A Task represents something to do before an arbitrary due date.
type Task struct {
	title      string // the title of the task (0 < length < 256)
	desc       string // the description of the task (0 <= length <= 65535)
	isPeriodic bool   // indicates if the task is periodic
	priority   byte   // the priority of the task: Low(0), Medium(1), High(2)
	status     byte   // the status of the task: Todo(3), Doing(4), Done(5)
}

func NewTask(title, desc string, isPeriodic bool, priority, status byte) (*Task,
	error) {
	if err := checkTitleValidity(title); err != nil {
		return nil, err
	}
	if len(desc) > DescMaxLength {
		return nil, ErrDescTooLong
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

// Removes the first character of s if it is a slash '/'. If it is not, returns
// s
func removeFirstSlashIfPresent(s string) string {
	if s != "" && s[0] == '/' {
		return s[1:]
	}
	return s
}

// Saves this task on disk
func (t *Task) saveAt(path string) error {
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	path = path + t.Title()
	offset := 0
	data := make([]byte, t.Length())
	titleLen := len(t.Title())
	data[offset] = byte(titleLen)
	offset++
	copy(data[offset:offset+titleLen], t.Title())
	offset += titleLen
	descLen := uint16(len(t.Description()))
	data[offset] = byte(descLen >> 8)
	offset++
	data[offset] = byte(descLen)
	offset++
	copy(data[offset:offset+int(descLen)], t.Description())
	offset += int(descLen)
	if t.IsPeriodic() {
		data[offset] = 1
	} else {
		data[offset] = 0
	}
	offset++
	data[offset] = byte(t.Priority())
	offset++
	data[offset] = byte(t.Status())
	return os.WriteFile(filepath.Clean(path), data, 0644)
}

// Saves on disk this task. TasksPath must be set before the call. Returns an
// error if something wrong happened
func (t *Task) SaveOnDisk() error {
	return t.saveAt(TasksPath)
}

// Returns the length in bytes needed to store this task
func (t *Task) Length() int {
	return 1 + len(t.Title()) + 2 + len(t.Description()) + 3
}
