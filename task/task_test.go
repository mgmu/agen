package task

import (
	"os"
	"strings"
	"testing"
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
		t.Fatalf("got \"%d\", want \"%d\"", ts.Priority(), Low)
	}
}

func TestNewTaskWithMediumPriorityHasMediumPriority(t *testing.T) {
	ts, err := NewTask("new", "", false, Medium, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != Medium {
		t.Fatalf("got \"%d\", want \"%d\"", ts.Priority(), Medium)
	}
}

func TestNewTaskWithHighPriorityHasHighPriority(t *testing.T) {
	ts, err := NewTask("new", "", false, High, Todo)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Priority() != High {
		t.Fatalf("got \"%d\", want \"%d\"", ts.Priority(), High)
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
		t.Fatalf("got \"%d\", want \"%d\"", ts.Status(), Todo)
	}
}

func TestNewTaskWithDoingStatusHasDoingStatus(t *testing.T) {
	ts, err := NewTask("new", "", false, Medium, Doing)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Doing {
		t.Fatalf("got \"%d\", want \"%d\"", ts.Status(), Doing)
	}
}

func TestNewTaskWithDoneStatusHasDoneStatus(t *testing.T) {
	ts, err := NewTask("new", "", false, High, Done)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if ts.Status() != Done {
		t.Fatalf("got \"%d\", want \"%d\"", ts.Status(), Done)
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
		t.Fatalf("got \"%d\", want \"%d\"", ts.Priority(), Medium)
	}
	if ts.Status() != Todo {
		t.Fatalf("got \"%d\", want \"%d\"", ts.Status(), Todo)
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
		t.Fatalf("got \"%d\", want \"%d\"", l, rl)
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
		t.Fatalf("got \"%d\", want \"%d\"", l, rl)
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
		t.Fatalf("got isPeriodic \"%t\", want \"%t\"", isPeriodic,
			task.IsPeriodic())
	}
	savedPriority := data[int(descLen)+1+titleLen+3]
	if task.Priority() != int(savedPriority) {
		t.Fatalf("got priority \"%d\", want \"%d\"", savedPriority,
			task.Priority())
	}
	savedStatus := data[int(descLen)+1+titleLen+4]
	if task.Status() != int(savedStatus) {
		t.Fatalf("got status \"%d\", want \"%d\"", savedStatus,	task.Status())
	}
}
