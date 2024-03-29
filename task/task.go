package task

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

	ErrTitleTooShort       = errors.New("title too short (min 1)")
	ErrTitleTooLong        = errors.New("title too long (max 255)")
	ErrInvalidPriority     = errors.New("priority must be Low, Medium or High")
	ErrInvalidStatus       = errors.New("status must be Todo, Doing or Done")
	ErrDescTooLong         = errors.New("description too long (max 65535)")
	ErrInvalidLoadPath     = errors.New("invalid load path")
	ErrInvalidTaskFileSize = errors.New("invalid task file size")
)

// A Task represents something to do before an arbitrary due date.
type Task struct {
	title      string // the title of the task (0 < length < 256)
	desc       string // the description of the task (0 <= length <= 65535)
	isPeriodic bool   // indicates if the task is periodic
	priority   byte   // the priority of the task: Low(0), Medium(1), High(2)
	status     byte   // the status of the task: Todo(3), Doing(4), Done(5)
	uuid       string // the uuid of the task
}

func NewTask(title, desc string, isPeriodic bool, priority, status byte) (*Task,
	error) {
	if err := checkTitleValidity(title); err != nil {
		return nil, err
	}
	if !isValidDescription(desc) {
		return nil, ErrDescTooLong
	}
	if !isValidPriority(priority) {
		return nil, ErrInvalidPriority
	}
	if !isValidStatus(status) {
		return nil, ErrInvalidStatus
	}
	new := Task{
		title:      title,
		desc:       desc,
		isPeriodic: isPeriodic,
		priority:   priority,
		status:     status,
		uuid:       uuid.NewString(),
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
		uuid:     uuid.NewString(),
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

// Returns the UUID of the task.
func (t *Task) Uuid() string {
	return t.uuid
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
// s.
func removeFirstSlashIfPresent(s string) string {
	if s != "" && s[0] == '/' {
		return s[1:]
	}
	return s
}

// Saves this task on disk
func (t *Task) saveAt(path string) error {
	if path == "" {
		return ErrInvalidLoadPath
	}
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	path = path + t.uuid
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
	offset++
	uuidLen := len(t.uuid)
	data[offset] = byte(uuidLen)
	offset++
	copy(data[offset:offset+uuidLen], t.uuid)
	return os.WriteFile(filepath.Clean(path), data, 0644)
}

// Saves on disk this task. TasksPath must be set before the call. Returns an
// error if something wrong happened
func (t *Task) SaveOnDisk() error {
	return t.saveAt(TasksPath)
}

// Returns the number of bytes needed to store this task.
func (t *Task) Length() int {
	return 1 + len(t.Title()) + 2 + len(t.Description()) + 4 + len(t.uuid)
}

// Loads to memory the task denoted by the given filepath and returns a pointer
// to it.
func loadTaskFrom(path string) (*Task, error) {
	if path == "" {
		return nil, ErrInvalidLoadPath
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, ErrInvalidLoadPath
	}
	size := info.Size()
	if size < 1 {
		return nil, ErrInvalidTaskFileSize
	}
	data := make([]byte, size)
	n, err := f.Read(data)
	if err != nil {
		return nil, err
	}
	data = data[:n]
	if len(data) < 1 {
		return nil, ErrInvalidTaskFileSize
	}
	offset := 0
	titleLen := int(data[offset])
	offset++
	if len(data) < offset+titleLen {
		return nil, ErrInvalidTaskFileSize
	}
	title := string(data[offset : offset+titleLen])
	offset += titleLen
	if len(data) < offset+2 {
		return nil, ErrInvalidTaskFileSize
	}
	descLen := uint16(data[offset])<<8 + uint16(data[offset+1])
	offset += 2
	if len(data) < offset+int(descLen) {
		return nil, ErrInvalidTaskFileSize
	}
	desc := string(data[offset : offset+int(descLen)])
	offset += int(descLen)
	if len(data) < offset+3 {
		return nil, ErrInvalidTaskFileSize
	}
	isPeriodic := (data[offset] == 1)
	offset++
	priority := data[offset]
	offset++
	status := data[offset]
	offset++
	if len(data) < offset+2 {
		return nil, ErrInvalidTaskFileSize
	}
	uuidLen := int(data[offset])
	offset++
	if len(data) < offset+uuidLen {
		return nil, ErrInvalidTaskFileSize
	}
	uuid := string(data[offset : offset+uuidLen])
	newTask, err := NewTask(title, desc, isPeriodic, priority, status)
	if err != nil {
		return nil, err
	}
	newTask.uuid = uuid
	return newTask, nil
}

// Loads into a slice of Task pointers the tasks saved at the given path
func loadTasksFrom(path string) ([]*Task, error) {
	if path == "" {
		return nil, ErrInvalidLoadPath
	}
	dir, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	info, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrInvalidLoadPath
	}
	entries, err := dir.ReadDir(0)
	if err != nil {
		return nil, err
	}
	var tasks []*Task
	for _, entry := range entries {
		ts, err := loadTaskFrom(filepath.Join(path, entry.Name()))
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, ts)
	}
	return tasks, nil
}

// Loads into a slice of Task pointers the tasks saved on disk
func LoadTasks() ([]*Task, error) {
	return loadTasksFrom(TasksPath)
}

// Counts the number of files that have the given prefix at the given directory
// path. Returns that number and a nil error or -1 and a non-nil error if an
// error occured.
func countFilesWithPrefixAt(path, prefix string) (int, error) {
	error_int := -1
	if prefix == "" || path == "" {
		return error_int, ErrInvalidLoadPath
	}
	dir, err := os.Open(path)
	if err != nil {
		return error_int, err
	}
	info, err := dir.Stat()
	if err != nil {
		return error_int, err
	}
	if !info.IsDir() {
		return error_int, ErrInvalidLoadPath
	}
	entries, err := dir.ReadDir(0)
	if err != nil {
		return -1, err
	}
	count := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prefix) {
			count++
		}
	}
	return count, nil
}

// Returns true if a task of given uuid of part of it already exists at the
// given directory path. Note that if there are multiple files that start with
// the given prefix, it returns also true but does not detect the fact that
// there is more that one file with that prefix.
func existsAt(path, uuid string) (bool, error) {
	res, err := countFilesWithPrefixAt(path, uuid)
	if err != nil {
		return false, err
	}
	return res > 0, nil
}

// Returns true if a task of given uuid or part of it already exists on disk.
func Exists(uuid string) (bool, error) {
	return existsAt(TasksPath, uuid)
}

// Returns true if a task of given uuid or has given uuid as prefix exists and
// that task is the only one that has the given uuid as a prefix.
func ExistsAndIsUnique(uuid string) (bool, error) {
	res, err := countFilesWithPrefixAt(TasksPath, uuid)
	if err != nil {
		return false, err
	}
	return res == 1, nil
}

// Returns a string that displays the title, the status and the priority of this
// task
func (t *Task) Display() string {
	prioDisp := ""
	switch t.Priority() {
	case Low:
		prioDisp = "low"
	case Medium:
		prioDisp = "medium"
	default:
		prioDisp = "high"
	}
	statusDisp := ""
	switch t.Status() {
	case Todo:
		statusDisp = "To do"
	case Doing:
		statusDisp = "Doing"
	default:
		statusDisp = "Done"
	}
	return fmt.Sprintf("[%s] %s <%s> %s", statusDisp, t.Title(), prioDisp,
		t.Uuid())
}

// Sets the description of this task to the given description. If the
// description is too long, an error is returned
func (t *Task) SetDescription(newDesc string) error {
	if len(newDesc) > DescMaxLength {
		return ErrDescTooLong
	}
	t.desc = newDesc
	return nil
}

// Sets the periodicity of this task
func (t *Task) SetPeriodicity(newPeriod bool) {
	t.isPeriodic = newPeriod
}

// Returns true if the given description is valid
func isValidDescription(desc string) bool {
	return len(desc) < 65536
}

// Returns true if the given priority is valid
func isValidPriority(prio byte) bool {
	return prio < 3
}

// Returns true if the given status is valid
func isValidStatus(status byte) bool {
	return status < 6 && status > 2
}

// Sets the new priority of this task. Returns an error if the priority is not
// valid
func (t *Task) SetPriority(newPrio byte) error {
	if !isValidPriority(newPrio) {
		return ErrInvalidPriority
	}
	t.priority = newPrio
	return nil
}

// Sets the new status of this task. Returns an error if the status is not valid
func (t *Task) SetStatus(newStatus byte) error {
	if !isValidStatus(newStatus) {
		return ErrInvalidStatus
	}
	t.status = newStatus
	return nil
}

// Parses the status denoted by the given string. If status is different than
// "todo", "doing" or "done", returns an error
func ParseStatus(status string) (byte, error) {
	switch status {
	case "todo":
		return Todo, nil
	case "doing":
		return Doing, nil
	case "done":
		return Done, nil
	default:
		return 0, errors.New("not a valid status string")
	}
}

// Loads the task of given uuid or part of it. If the task does not exist or
// something happens during the load, returns an error.
func LoadTask(uuid string) (*Task, error) {
	uuidLen := len(uuid)
	if uuidLen == 0 || uuidLen > 36 {
		return nil, errors.New("invalid uuid length")
	}
	if uuidLen < 36 {
		return loadTaskFromPartialUuid(uuid)
	}
	return loadTaskFrom(filepath.Join(TasksPath, uuid))
}

// Parses the priority denoted by the given string. If priority is different
// than "low", "medium" or "high", returns an error
func ParsePriority(priority string) (byte, error) {
	switch priority {
	case "low":
		return Low, nil
	case "medium":
		return Medium, nil
	case "high":
		return High, nil
	default:
		return 0, errors.New("not a valid priority string")
	}
}

// Removes the file of given name at the given path if it exists, otherwise
// returns an error.
func removeAt(path, name string) error {
	exists, err := Exists(name)
	if err != nil {
		return err
	}
	if exists {
		return os.Remove(filepath.Join(path, name))
	}
	return nil
}

// Searchs for the full uuid if it is partial and removes the corresponding
// task. If an error occurs, it returns the error.
func removeFromPartialUuid(uuid string) error {
	taskDir, err := os.ReadDir(TasksPath)
	if err != nil {
		return err
	}
	for _, entry := range taskDir {
		filename := entry.Name()
		if strings.HasPrefix(filename, uuid) {
			return removeAt(TasksPath, filename)
		}
	}
	return errors.New("task not found")
}

// Removes the task of given uuid or part of if. If multiple tasks have the
// given uuid as prefix, no tasks are removed and an error is returned.
func Remove(uuid string) error {
	existsAndUnique, err := ExistsAndIsUnique(uuid)
	if err != nil {
		return err
	}
	if !existsAndUnique {
		return errors.New("uuid prefix not unique")
	}
	uuidLen := len(uuid)
	if uuidLen == 0 || uuidLen > 36 {
		return errors.New("invalid uuid length")
	}
	if uuidLen < 36 {
		return removeFromPartialUuid(uuid)
	}
	return removeAt(TasksPath, uuid)
}

// Parses the strings and returns the slice of status marks found, without
// duplicates
func ParseStatusFrom(strings []string) ([]byte, error) {
	var res []byte
	for _, str := range strings {
		if IsValidStatus(str) {
			stat, err := ParseStatus(str)
			if err != nil {
				return nil, err
			}
			if !slices.Contains(res, stat) {
				res = append(res, stat)
			}
		}
	}
	return res, nil
}

// Parses the strings and returns the slice of priority marks found, without
// duplicates
func ParsePriorityFrom(strings []string) ([]byte, error) {
	var res []byte
	for _, str := range strings {
		if IsValidPriority(str) {
			prio, err := ParsePriority(str)
			if err != nil {
				return nil, err
			}
			if !slices.Contains(res, prio) {
				res = append(res, prio)
			}
		}
	}
	return res, nil
}

// Returns true if status equals one of "todo", "doing" or "done"
func IsValidStatus(status string) bool {
	return status == "todo" || status == "doing" || status == "done"
}

// Returns true if priority equals one of "low", "medium" or "high"
func IsValidPriority(priority string) bool {
	return priority == "low" || priority == "medium" || priority == "high"
}

// Filters the given tasks and returns the remaining tasks
func FilterTasks(tasks []*Task, filters []string) ([]*Task, error) {
	sFilters, err := ParseStatusFrom(filters)
	if err != nil {
		return nil, err
	}
	pFilters, err := ParsePriorityFrom(filters)
	if err != nil {
		return nil, err
	}
	if len(sFilters) == 0 && len(pFilters) == 0 {
		return tasks, nil
	}
	var res []*Task
	for _, task := range tasks {
		if len(sFilters) != 0 && !slices.Contains(sFilters, task.status) {
			continue
		}
		if len(pFilters) != 0 && !slices.Contains(pFilters, task.priority) {
			continue
		}
		res = append(res, task)
	}
	return res, nil
}

func loadTaskFromPartialUuid(uuid string) (*Task, error) {
	taskDir, err := os.ReadDir(TasksPath)
	if err != nil {
		return nil, err
	}
	for _, entry := range taskDir {
		filename := entry.Name()
		if strings.HasPrefix(filename, uuid) {
			path := filepath.Join(TasksPath, filename)
			return loadTaskFrom(path)
		}
	}
	return nil, errors.New("task not found")
}
