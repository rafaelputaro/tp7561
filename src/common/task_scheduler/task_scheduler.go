package task_scheduler

import (
	"errors"
	"sync"
)

const MAX_TASK = 1000
const MSG_TASK_SCHEDULER_BUSY_OR_CLOSED = "no more tasks can be accepted"

// Se encarga de mantener una serie de tareas a ser descargadas de un canal para ser ejecutadas
type TaskScheduler struct {
	tasksChan        chan func()
	taggedTasks      map[string]bool
	mutexTaggedTasks sync.Mutex
}

// Retorna una instancia de Task Scheduler lista para ser utilizada
func NewTaskScheduler() *TaskScheduler {
	scheduler := TaskScheduler{
		tasksChan:   make(chan func(), MAX_TASK),
		taggedTasks: map[string]bool{},
	}
	go func() {
		notClosed := true
		for notClosed {
			task, ok := <-scheduler.tasksChan
			if ok {
				task()
			}
			notClosed = ok
		}
	}()
	return &scheduler
}

// Cierra el canal descartando las tareas pendientes
func (scheduler *TaskScheduler) DisposeTaskScheduler() {
	close(scheduler.tasksChan)
}

// Agrega una tarea a ser ejecutada
func (scheduler *TaskScheduler) AddTask(task func()) (err error) {
	// Recuperación de panic ante canal cerrado
	defer func() {
		if recover() != nil {
			err = errors.New(MSG_TASK_SCHEDULER_BUSY_OR_CLOSED) // Canal lleno o cerrado, no bloquea
		} else {
			err = nil
		}
	}()
	select {
	case scheduler.tasksChan <- task:
		return nil // Envío exitoso
	default:
		return errors.New(MSG_TASK_SCHEDULER_BUSY_OR_CLOSED) // Canal lleno o cerrado, no bloquea
	}
}

// Agrega una tarea a ser ejecutada
func (scheduler *TaskScheduler) AddTaggedTask(task func(), tag string) (err error) {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	scheduler.taggedTasks[tag] = true
	return scheduler.AddTask(task)
}

// Remueve una etiqueta de la lista de tareas etiquetadas
func (scheduler *TaskScheduler) RemoveTaggedTask(tag string) {
	delete(scheduler.taggedTasks, tag)
}

// Retorna verdadero si encuentra el tag
func (scheduler *TaskScheduler) HasTag(tag string) bool {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	_, ok := scheduler.taggedTasks[tag]
	return ok
}

// Retorna verdadero si no hay tareas pendientes
func (scheduler *TaskScheduler) Empty() bool {
	return len(scheduler.tasksChan) == 0
}
