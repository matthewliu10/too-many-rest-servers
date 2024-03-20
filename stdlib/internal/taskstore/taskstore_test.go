package taskstore

import (
	"testing"
	"time"
	"reflect"
)

func TestCreateAndGet(t *testing.T) {
	ts := New()

	text, tags, due := "testing123", []string{"test"}, time.Now()
	id0 := ts.CreateTask(text, tags, due)

	// test ts.GetTask
	task, err := ts.GetTask(id0)
	if err != nil {
		t.Fatal(err.Error())
	}

	if task.Id != id0 || task.Text != text || !reflect.DeepEqual(tags, task.Tags)|| task.Due != due {
		t.Errorf("\nexpect: id=%d, text=%s, tags=%#v, due=%v\nactual: id=%d, text=%s, tags=%#v, due=%v",
			id0,
			text,
			tags,
			due,
			task.Id,
			task.Text,
			task.Tags,
			task.Due,
		)
	}

	// test ts.GetTask with an invalid task
	task, err = ts.GetTask(id0 + 1)
	if err == nil {
		t.Errorf("\nexpect: ts.GetTask return an error with invalid task\nactual: task=%+v", task)
	}

	// test ts.GetAllTasks
	ts.CreateTask(text, tags, due)
	tasks := ts.GetAllTasks()
	if len(tasks) != 2 {
		t.Errorf("\nexpect: len(tasks)=%d\nactual: len(tasks)=%d", 2, len(tasks))
	}
}

func TestDelete(t *testing.T) {
	ts := New()

	text, tags, due := "testing123", []string{"test"}, time.Now()
	id0 := ts.CreateTask(text, tags, due)
	ts.CreateTask(text, tags, due)
	ts.CreateTask(text, tags, due)

	// test DeleteTask
	ts.DeleteTask(id0)
	_, err := ts.GetTask(id0)
	if err == nil {
		t.Error("\nexpect: fail to get deleted task\nactual: getTask ran successfully on deleted task")
	}

	// test DeleteAll
	ts.DeleteAllTasks()
	tasks := ts.GetAllTasks()
	if len(tasks) != 0 {
		t.Errorf("\nexpect: len(tasks)=0 after deleting all tasks\nactual: len(tasks)=%d", len(tasks))
	}
}

func TestTags(t *testing.T) {
	ts := New()

	text, tags1, tags2, tags3, due := "testing123", []string{"a", "b"}, []string{"c", "a"}, []string{"e"}, time.Now()
	ts.CreateTask(text, tags1, due)
	ts.CreateTask(text, tags2, due)
	ts.CreateTask(text, tags3, due)

	// test GetTasksByTag
	tasks := ts.GetTasksByTag("a")
	if len(tasks) != 2 {
		t.Errorf("\nexpect: len(tasks)=%d\nactual: len(tasks)=%d", 2, len(tasks))
	}
}

func TestDueDate(t *testing.T) {
	ts := New()

	text, tags, due, due2 := "testing123", []string{"test"}, time.Now(), time.Now().AddDate(0, 0, 1)
	id0 := ts.CreateTask(text, tags, due)
	id1 := ts.CreateTask(text, tags, due)
	ts.CreateTask(text, tags, due2)

	// test GetByDueDate
	tasks := ts.GetByDueDate(due.Year(), due.Month(), due.Day())
	if len(tasks) != 2 || tasks[0].Id != id0 || tasks[1].Id != id1 {
		t.Errorf("\nexpect: len(tasks)=%d, tasks[0].Id=%d, tasks[1].Id=%d\nactual: len(tasks)=%d, tasks[0].Id=%d, tasks[1].Id=%d",
			2,
			id0,
			id1,
			len(tasks),
			tasks[0].Id,
			tasks[1].Id,
		)
	}
}
