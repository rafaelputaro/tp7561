package task_scheduler

import (
	"errors"
	"sync"
	"tp/common"
)

const MAX_TASK = 10000
const MAX_TASK_RETRIES = 200
const MSG_TASK_SCHEDULER_BUSY_OR_CLOSED = "no more tasks can be accepted"
const MSG_TAG_EXISTS = "tag exists: %v"
const MSG_TASK_ADDED = "added task: %v"
const MSG_RETRY_TASK = "retry task: tag %v | remaining retries: %v"
const MSG_RULE_OUT_TASK = "rule out task: tag %v | remaining retries: %v"

// Retorna el tag de la función y si la misma es candidata a ser reintentada
type TaskFunc func() (string, bool)

// Se encarga de mantener una serie de tareas a ser descargadas de un canal para ser ejecutadas
type TaskScheduler struct {
	tasksChan        chan TaskFunc
	taggedTasks      map[string]int
	mutexTaggedTasks *sync.Mutex
}

// Retorna una instancia de Task Scheduler lista para ser utilizada
func NewTaskScheduler() *TaskScheduler {
	scheduler := TaskScheduler{
		tasksChan:        make(chan TaskFunc, MAX_TASK),
		taggedTasks:      map[string]int{},
		mutexTaggedTasks: &sync.Mutex{},
	}
	go func() {
		notClosed := true
		for notClosed {
			task, ok := <-scheduler.tasksChan
			if ok {
				tag, retry := task()
				scheduler.checkRetryTask(retry, task, tag)
			}
			notClosed = ok
			//common.SleepBetweenTasks()
		}
	}()
	return &scheduler
}

// Chequea si la tarea debe ser reintentada
func (scheduler *TaskScheduler) checkRetryTask(retry bool, task TaskFunc, tag string) {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	if !retry {
		scheduler.doRemoveTaggedTask(tag)
		return
	}
	if scheduler.taggedTasks[tag] > 0 {
		scheduler.taggedTasks[tag] = scheduler.taggedTasks[tag] - 1
		scheduler.addTask(task)
		common.Log.Debugf(MSG_RETRY_TASK, tag, scheduler.taggedTasks[tag])
		return
	}
	common.Log.Debugf(MSG_RULE_OUT_TASK, tag, scheduler.taggedTasks[tag])
	scheduler.doRemoveTaggedTask(tag)
}

// Cierra el canal descartando las tareas pendientes
func (scheduler *TaskScheduler) DisposeTaskScheduler() {
	close(scheduler.tasksChan)
}

// Agrega una tarea a ser ejecutada
func (scheduler *TaskScheduler) addTask(task TaskFunc) (err error) {
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
func (scheduler *TaskScheduler) AddTask(task TaskFunc, tag string) error {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	if scheduler.doHasTag(tag) {
		common.Log.Debugf(MSG_TAG_EXISTS, tag)
		return errors.New(MSG_TAG_EXISTS)
	}
	scheduler.taggedTasks[tag] = MAX_TASK_RETRIES
	common.Log.Debugf(MSG_TASK_ADDED, tag)
	err := scheduler.addTask(task)
	return err
}

// Remueve una etiqueta de la lista de tareas etiquetadas
func (scheduler *TaskScheduler) removeTask(tag string) {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	delete(scheduler.taggedTasks, tag)
}

// Remueve una etiqueta de la lista de tareas etiquetadas
func (scheduler *TaskScheduler) doRemoveTaggedTask(tag string) {
	delete(scheduler.taggedTasks, tag)
}

// Retorna verdadero si encuentra el tag
func (scheduler *TaskScheduler) doHasTag(tag string) bool {
	_, ok := scheduler.taggedTasks[tag]
	return ok
}

// Retorna verdadero si encuentra el tag
func (scheduler *TaskScheduler) HasTag(tag string) bool {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	return scheduler.doHasTag(tag)
}

// Retorna verdadero si no hay tareas pendientes
func (scheduler *TaskScheduler) Empty() bool {
	return len(scheduler.tasksChan) == 0
}
