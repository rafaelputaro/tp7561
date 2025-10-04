package task_scheduler

import (
	"errors"
	"sync"
)

const MAX_TASK = 1000
const MSG_TASK_SCHEDULER_BUSY = "no more tasks can be accepted"

type TaskScheduler struct {
	tasksChan chan func()
	mutex     sync.Mutex
}

func NewTaskScheduler() *TaskScheduler {
	scheduler := TaskScheduler{
		tasksChan: make(chan func(), MAX_TASK),
	}
	go func() {
		for {

		}
	}()
	return &scheduler
}

func (scheduler *TaskScheduler) AddTask(task func()) error {
	scheduler.mutex.Lock()
	defer scheduler.mutex.Unlock()
	if len(scheduler.tasksChan) < MAX_TASK {
		scheduler.tasksChan <- task
		return nil
	}
	return errors.New(MSG_TASK_SCHEDULER_BUSY)
}
