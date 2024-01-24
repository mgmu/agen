package task

import (
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
	exp := "title too short"
	if !equal(exp, err.Error()) {
		t.Fatalf("got \"%s\", want \"%s\"", err.Error(), exp)
	}
}

func TestTaskWithTitleTooLongReturnsError(t *testing.T) {
	longTitle := `tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt
tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt
tttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttttt
tttttttttttttttttttttttttttttttttt`
	_, err := NewTask(longTitle, "", false, Low, Todo)
	if err == nil {
		t.Fatalf("expected error")
	}
	exp := "title too long"
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
