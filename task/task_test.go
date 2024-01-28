package task

import (
	"os"
	"strings"
	"testing"
	"path/filepath"
)

// Returns true if s1 and s2 are exactly the same
func equal(s1, s2 string) bool {
	return strings.Compare(s1, s2) == 0
}

func TestTaskWithTitleTooShortReturnsError(t *testing.T) {
	_, err := NewTask("", "", false, Low, Todo)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "title too short (min 1)"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskWithTitleTooLongReturnsError(t *testing.T) {
	longTitle := strings.Repeat("t", TitleMaxLength + 1)
	_, err := NewTask(longTitle, "", false, Low, Todo)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "title too long (max 255)"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskHasGivenName(t *testing.T) {
	_, err := NewTask("new task", "", false, Low, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestNewTaskHasGivenDescription(t *testing.T) {
	ts, err := NewTask("new task", "some description", false, Low, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !equal(ts.Description(), "some description") {
		t.Fatalf("got \"%s\", want \"%s\"", ts.Description(),
			"some description")
	}
}

func TestInvalidPriorityForNewTaskReturnsError(t *testing.T) {
	_, err := NewTask("new", "", false, 3, Todo)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "priority must be Low, Medium or High"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestNewTaskWithLowPriorityHasLowPriority(t *testing.T) {
	ts, err := NewTask("new", "", false, Low, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != Low {
		t.Fatalf("got %d, want %d", ts.Priority(), Low)
	}
}

func TestNewTaskWithMediumPriorityHasMediumPriority(t *testing.T) {
	ts, err := NewTask("new", "", false, Medium, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != Medium {
		t.Fatalf("got %d, want %d", ts.Priority(), Medium)
	}
}

func TestNewTaskWithHighPriorityHasHighPriority(t *testing.T) {
	ts, err := NewTask("new", "", false, High, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != High {
		t.Fatalf("got %d, want %d", ts.Priority(), High)
	}
}

func TestNewTaskWithStatus6ReturnsError(t *testing.T) {
	_, err := NewTask("new", "", false, High, 6)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "status must be Todo, Doing or Done"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestNewTaskWithStatus2ReturnsError(t *testing.T) {
	_, err := NewTask("new", "", false, High, 2)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "status must be Todo, Doing or Done"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestNewTaskWithTodoStatusHasTodoStatus(t *testing.T) {
	ts, err := NewTask("new", "", false, Low, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Todo {
		t.Fatalf("got %d, want %d", ts.Status(), Todo)
	}
}

func TestNewTaskWithDoingStatusHasDoingStatus(t *testing.T) {
	ts, err := NewTask("new", "", false, Medium, Doing)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Doing {
		t.Fatalf("got %d, want %d", ts.Status(), Doing)
	}
}

func TestNewTaskWithDoneStatusHasDoneStatus(t *testing.T) {
	ts, err := NewTask("new", "", false, High, Done)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Done {
		t.Fatalf("got %d, want %d", ts.Status(), Done)
	}
}

func TestNewDefaultTaskHasDefaultParametersAndGivenTitle(t *testing.T) {
	ts, err := NewDefault("default")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !equal(ts.Title(), "default") {
		t.Fatalf("got \"%s\", want \"%s\"", ts.Title(), "default")
	}
	if !equal(ts.Description(), "") {
		t.Fatalf("got \"%s\", want \"%s\"", ts.Description(), "")
	}
	if ts.IsPeriodic() {
		t.Fatalf("got \"%t\", want \"%t\"", ts.IsPeriodic(), false)
	}
	if ts.Priority() != Medium {
		t.Fatalf("got %d, want %d", ts.Priority(), Medium)
	}
	if ts.Status() != Todo {
		t.Fatalf("got %d, want %d", ts.Status(), Todo)
	}
}

func TestNewDefaultTaskWithNameInferiorTo1CharacterReturnsError(t *testing.T) {
	_, err := NewDefault("")
	if err == nil {
		t.Fatalf("error expected")
	}
	exp := "title too short (min 1)"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestNewDefaultTaskWithTitleTooLongReturnsError(t *testing.T) {
	longTitle := strings.Repeat("t", TitleMaxLength + 1)
	_, err := NewDefault(longTitle)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "title too long (max 255)"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestRemoveFirstSlashReturnsEmptyStringIfEmptyStringGiven(t *testing.T) {
	s := removeFirstSlashIfPresent("")
	if s != "" {
		t.Fatalf("got \"%s\", want \"%s\"", s, "")
	}
}

func TestRemoveFirstSlashWithoutFirstSlashReturnsSameString(t *testing.T) {
	s := removeFirstSlashIfPresent("no slash")
	if s != "no slash" {
		t.Fatalf("got \"%s\", want \"%s\"", s, "no slash")
	}
}

func TestRemoveFirstSlashReturnsSameStringWithoutFirstSlash(t *testing.T) {
	s := removeFirstSlashIfPresent("/slash")
	if s != "slash" {
		t.Fatalf("got \"%s\", want \"%s\"", s, "slash")
	}
}

func TestNewTaskWithDescriptionTooLongReturnsError(t *testing.T) {
	desc := strings.Repeat("a", DescMaxLength + 1)
	_, err := NewTask("title", desc, false, Medium, Todo)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "description too long (max 65535)"
	if err.Error() != exp {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestSaveAtSavesTheTaskAtPath(t *testing.T) {
	task, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = task.saveAt("/tmp")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.Remove("/tmp/test")
	_, err = os.ReadFile("/tmp/test")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLengthNewDefaultTask(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	l := ts.Length()
	rl := 1 + 4 + 2 + 0 + 1 + 1 + 1
	if l != rl {
		t.Fatalf("got %d, want %d", l, rl)
	}
}

func TestLengthNewTask(t *testing.T) {
	ts, err := NewTask("test", "this is a long desc", false, Medium, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	l := ts.Length()
	rl := 1 + 4 + 2 + 19 + 1 + 1 + 1
	if l != rl {
		t.Fatalf("got %d, want %d", l, rl)
	}
}

func TestSaveAtSavesTheTaskAtPathAndTheTaskContent(t *testing.T) {
	task, err := NewTask("test", "this is a description", false, Medium, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = task.saveAt("/tmp")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.Remove("/tmp/test")
	data, err := os.ReadFile("/tmp/test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	titleLen := int(data[0])
	savedTitle := string(data[1:1+titleLen])
	if task.Title() != savedTitle {
		t.Fatalf("got title \"%s\", want \"%s\"", savedTitle, task.Title())
	}
	descLen := uint16(data[1+titleLen]) << 8 + uint16(data[1+titleLen+1])
	savedDesc := string(data[1+titleLen+2:int(descLen)+1+titleLen+2])
	if task.Description() != savedDesc {
		t.Fatalf("got description \"%s\", want \"%s\"", savedDesc,
			task.Description())
	}
	isPeriodic := (data[int(descLen)+1+titleLen+2] == 1)
	if task.IsPeriodic() != isPeriodic {
		t.Fatalf("got isPeriodic %t, want %t", isPeriodic,
			task.IsPeriodic())
	}
	savedPriority := data[int(descLen)+1+titleLen+3]
	if task.Priority() != int(savedPriority) {
		t.Fatalf("got priority %d, want %d", savedPriority,
			task.Priority())
	}
	savedStatus := data[int(descLen)+1+titleLen+4]
	if task.Status() != int(savedStatus) {
		t.Fatalf("got status %d, want %d", savedStatus,	task.Status())
	}
}

func TestLoadFromEmtpyDirReturnsEmptySlice(t *testing.T) {
	dirname, err := os.MkdirTemp("", "tasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	tasks, err := loadTasksFrom(dirname)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tasks) != 0 {
		t.Fatalf("got %v, want empty slice", tasks)
	}
}

func TestLoadFromDirWith2TasksReturnsSliceOfLength2(t *testing.T) {
	dirname, err := os.MkdirTemp("", "tasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	ts1, err := NewDefault("test1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ts2, err := NewDefault("test2")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err = ts1.saveAt(dirname); err != nil {
		t.Fatalf(err.Error())		
	}
	if err = ts2.saveAt(dirname); err != nil {
		t.Fatalf(err.Error())		
	}
	tasks, err := loadTasksFrom(dirname)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tasks) != 2 {
		t.Fatalf("got slice of length %d, want %d", len(tasks), 2)
	}
}

func TestLoadFromWithEmptyStringReturnsError(t *testing.T) {
	_, err := loadTasksFrom("")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestLoadTaskFromEmptyStringReturnsError(t *testing.T) {
	_, err := loadTaskFrom("")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestLoadTaskFromPathThatIssNotARegularFileReturnsError(t *testing.T) {
	_, err := loadTaskFrom(".")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestLoadTaskFromPathThatIsShortRegFileReturnsError(t *testing.T) {
	dirname, err := os.MkdirTemp("", "tempTasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	_, err = os.Create(filepath.Join(dirname, "test"))
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = loadTaskFrom(filepath.Join(dirname, "test"))
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid task file size"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestLoadTaskFromPathThatIsTaskReturnsCorrespondingTask(t *testing.T) {
	ts, err := NewTask("test", "its description", true, High, Done)
	if err != nil {
		t.Fatalf(err.Error())
	}
	dirname, err := os.MkdirTemp("", "tempTasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	if err = ts.saveAt(dirname); err != nil {
		t.Fatalf(err.Error())
	}
	tsptr, err := loadTaskFrom(filepath.Join(dirname, ts.Title()))
	if err != nil {
		t.Fatalf(err.Error())
	}
	if tsptr == nil {
		t.Fatalf("expected valid task pointer")
	}
	if tsptr.Title() != ts.Title() || tsptr.Description() != ts.Description() ||
		tsptr.IsPeriodic() != ts.IsPeriodic() ||
		tsptr.Priority() != ts.Priority() || tsptr.Status() != ts.Status() {
		t.Fatalf("expected same task")
	}
}

func TestLoadFromDirWith2TasksReturnsCorrespondingTasks(t *testing.T) {
	dirname, err := os.MkdirTemp("", "tasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	ts1, err := NewDefault("test1")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ts2, err := NewDefault("test2")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if err = ts1.saveAt(dirname); err != nil {
		t.Fatalf(err.Error())		
	}
	if err = ts2.saveAt(dirname); err != nil {
		t.Fatalf(err.Error())		
	}
	tasks, err := loadTasksFrom(dirname)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tasks) != 2 {
		t.Fatalf("got slice of length %d, want %d", len(tasks), 2)
	}
	if tasks[0].Title() != ts1.Title() ||
		tasks[0].Description() != ts1.Description() ||
		tasks[0].IsPeriodic() != ts1.IsPeriodic() ||
		tasks[0].Priority() != ts1.Priority() ||
		tasks[0].Status() != ts1.Status() {
		t.Fatalf("expected same task")
	}	
	if tasks[1].Title() != ts2.Title() ||
		tasks[1].Description() != ts2.Description() ||
		tasks[1].IsPeriodic() != ts2.IsPeriodic() ||
		tasks[1].Priority() != ts2.Priority() ||
		tasks[1].Status() != ts2.Status() {
		t.Fatalf("expected same task")
	}	
}

func TestTaskExistenceWithEmptyTitleStringReturnsError(t *testing.T) {
	_, err := existsAt("/tmp", "")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskExistenceWithEmptyPathStringReturnsError(t *testing.T) {
	_, err := existsAt("", "hi")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskExistenceWithPathThatIsNotADirReturnsError(t *testing.T) {
	dirname, err := os.MkdirTemp("", "test-agen")
	if err != nil {
		t.Fatalf(err.Error())
	}
	defer os.RemoveAll(dirname)
	path := filepath.Join(dirname, "hi")
	_, err = os.Create(path)
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = existsAt(path, "hi")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "invalid load path"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskExistenceOfNotExistentTaskReturnsFalse(t *testing.T) {
	dirname, err := os.MkdirTemp("", "tasks")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ok, err := existsAt(dirname, "not-present")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ok {
		t.Fatalf("got %t, want %t", ok, false)
	}
}

func TestSaveAtEmptyStringReturnsError(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	exp := "invalid load path"
	err = ts.saveAt("")
	if err == nil {
		t.Fatalf("expected error")
	}
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestListDisplayOfDefaultTaskShowsGoodTitleStatusAndPriority(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	display := ts.Display()
	exp := "[To do] test <medium>"
	if display != exp {
		t.Fatalf("got \"%s\", want \"%s\"", display, exp)		
	}
}

func TestSetDescriptionTooLongReturnsError(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = ts.SetDescription(strings.Repeat("a", 65536))
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "description too long (max 65535)"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestSetValidDescriptionReturnsNilErrorAndModifiesTask(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	newDesc := strings.Repeat("a", 65535)
	err = ts.SetDescription(newDesc)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if newDesc != ts.Description() {
		t.Fatalf("got \"%s\", want \"%s\"", ts.Description(), newDesc)
	}
}

func TestSetPeriodicityUpdatesTask(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	ts.SetPeriodicity(true)
	if !ts.IsPeriodic() {
		t.Fatalf("got %t, want %t", ts.IsPeriodic(), true)
	}
}

func TestSetPriorityToInvalidPriorityReturnsError(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = ts.SetPriority(3)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "priority must be Low, Medium or High"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestSetPriorityToHighUpdatesPriorityOfTask(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = ts.SetPriority(High)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != High {
		t.Fatalf("got %d, want %d", ts.Priority(), High)
	}
}

func TestSetStatusToInvalidStatusReturnsError(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = ts.SetStatus(2)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "status must be Todo, Doing or Done"
	if exp != err.Error() {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestSetStatusToDoneUpdatesStatusOfTask(t *testing.T) {
	ts, err := NewDefault("test")
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = ts.SetStatus(Done)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Done {
		t.Fatalf("got %d, want %d", ts.Status(), Done)
	}
}

func TestParseStatusFromEmptyStringReturnsError(t *testing.T) {
	_, err := ParseStatus("")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "not a valid status string"
	if err.Error() != exp {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestParseStatusWithTodoStringReturnsTodo(t *testing.T) {
	s, err := ParseStatus("todo")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != Todo {
		t.Fatalf("got %d, want %d", s, Todo)
	}
}

func TestParseStatusWithDoingStringReturnsDoing(t *testing.T) {
	s, err := ParseStatus("doing")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != Doing {
		t.Fatalf("got %d, want %d", s, Doing)
	}
}

func TestParseStatusWithDoneStringReturnsDone(t *testing.T) {
	s, err := ParseStatus("done")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != Done {
		t.Fatalf("got %d, want %d", s, Done)
	}
}

func TestParsePriorityFromEmptyStringReturnsError(t *testing.T) {
	_, err := ParsePriority("")
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "not a valid priority string"
	if err.Error() != exp {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestParsePriorityWithLowStringReturnsLow(t *testing.T) {
	s, err := ParsePriority("low")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != Low {
		t.Fatalf("got %d, want %d", s, Low)
	}
}

func TestParsePriorityWithMediumStringReturnsMedium(t *testing.T) {
	s, err := ParsePriority("medium")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != Medium {
		t.Fatalf("got %d, want %d", s, Medium)
	}
}

func TestParsePriorityWithHighStringReturnsHigh(t *testing.T) {
	s, err := ParsePriority("high")
	if err != nil {
		t.Fatalf(err.Error())
	}
	if s != High {
		t.Fatalf("got %d, want %d", s, High)
	}
}
