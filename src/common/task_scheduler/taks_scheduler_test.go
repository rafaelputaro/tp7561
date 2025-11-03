package task_scheduler

import (
	"strconv"
	"testing"
	"time"
)

var common_var = ""

func createTask(msg string) func() (string, bool) {
	return func() (string, bool) {
		common_var += msg
		return msg, true
	}
}

func TestScheduler(t *testing.T) {
	check := common_var
	scheduler := NewTaskScheduler()

	for i := range 20 {
		iS := strconv.Itoa(i)
		check += iS
		if err := scheduler.addTask(createTask(iS)); err != nil {
			t.Errorf("Error on add task: %v", err)
		}
	}
	for !scheduler.Empty() {
		time.Sleep(time.Second * 1)
	}
	if check != common_var {
		t.Errorf("Not Match: Expected: %v | Found: %v", common_var, check)
	}
	scheduler.DisposeTaskScheduler()
	if scheduler.addTask(createTask("hola")) == nil {
		t.Errorf("Add task on channel closed")
	}
}
