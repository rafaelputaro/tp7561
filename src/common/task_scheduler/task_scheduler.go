package task_scheduler

import (
	"errors"
	"sync"
	"time"
	"tp/common"
)

const MAX_TASK = 20000
const MAX_TASK_RETRIES = 200
const MSG_TASK_SCHEDULER_BUSY_OR_CLOSED = "no more tasks can be accepted"
const MSG_TAG_EXISTS = "tag exists: %v"
const MSG_TASK_ADDED = "added task: %v"
const MSG_RETRY_TASK = "retry task: tag %v | remaining retries: %v"
const MSG_RULE_OUT_TASK = "rule out task: tag %v | remaining retries: %v"
const MSG_RULE_OUT_TASK_TIME_EXPIRATION = "rule out task due to time expiration: tag %v | remaining retries: %v"

// Retorna el tag de la función y si la misma es candidata a ser reintentada
type TaskFunc func() (string, bool)

type TaskData struct {
	taskFunc   TaskFunc
	expiration *time.Time
	mandatory  bool
}

// Se encarga de mantener una serie de tareas a ser descargadas de un canal para ser ejecutadas
type TaskScheduler struct {
	tasksChan        chan TaskData
	taggedTasks      map[string]int
	mutexTaggedTasks *sync.Mutex
}

// Retorna una nueva instancia de tarea sin tiempo de expiración
func newTaskDataWithoutExpirationTime(task TaskFunc, mandatory bool) *TaskData {
	return &TaskData{
		taskFunc:   task,
		expiration: nil,
		mandatory:  mandatory,
	}
}

// Retorna una nueva tarea con tiempo de expiración siendo el parámetro "deltatime" en
// milisegundos la validez de la misma.
func newTaskDataWithExpirationTime(task TaskFunc, deltaTime float32) *TaskData {
	now := time.Now()
	expiration := now.Add(time.Millisecond * time.Duration(deltaTime))
	return &TaskData{
		taskFunc:   task,
		expiration: &expiration,
		mandatory:  false,
	}
}

// Retorna una instancia de Task Scheduler lista para ser utilizada
func NewTaskScheduler() *TaskScheduler {
	scheduler := TaskScheduler{
		tasksChan:        make(chan TaskData, MAX_TASK),
		taggedTasks:      map[string]int{},
		mutexTaggedTasks: &sync.Mutex{},
	}
	go func() {
		notClosed := true
		for notClosed {
			taskData, ok := <-scheduler.tasksChan
			task := taskData.taskFunc
			if ok {
				tag, retry := task()
				scheduler.checkRetryTask(retry, taskData, tag)
			}
			notClosed = ok
		}
	}()
	return &scheduler
}

// Chequea si la tarea debe ser reintentada
func (scheduler *TaskScheduler) checkRetryTask(retry bool, taskData TaskData, tag string) {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	if !retry {
		scheduler.doRemoveTask(tag)
		return
	}
	// chequear si es una tarea obligatoria
	if taskData.mandatory {
		scheduler.doAddTask(taskData)
		common.Log.Debugf(MSG_RETRY_TASK, tag, scheduler.taggedTasks[tag])
		return
	}
	// chequear vencimiento por tiempo
	if taskData.expiration != nil {
		now := time.Now()
		if taskData.expiration.Before(now) {
			common.Log.Debugf(MSG_RULE_OUT_TASK_TIME_EXPIRATION, tag, scheduler.taggedTasks[tag])
			scheduler.doRemoveTask(tag)
			return
		}
	}
	// chequear vencimiento por cantidad de intentos
	if scheduler.taggedTasks[tag] > 0 {
		scheduler.taggedTasks[tag] = scheduler.taggedTasks[tag] - 1
		scheduler.doAddTask(taskData)
		common.Log.Debugf(MSG_RETRY_TASK, tag, scheduler.taggedTasks[tag])
		return
	}
	common.Log.Debugf(MSG_RULE_OUT_TASK, tag, scheduler.taggedTasks[tag])
	scheduler.doRemoveTask(tag)
}

// Cierra el canal descartando las tareas pendientes
func (scheduler *TaskScheduler) DisposeTaskScheduler() {
	close(scheduler.tasksChan)
}

// Agrega una tarea a ser ejecutada
func (scheduler *TaskScheduler) doAddTask(taskData TaskData) (err error) {
	// Recuperación de panic ante canal cerrado
	defer func() {
		if recover() != nil {
			err = errors.New(MSG_TASK_SCHEDULER_BUSY_OR_CLOSED) // Canal lleno o cerrado, no bloquea
		} else {
			err = nil
		}
	}()
	select {
	case scheduler.tasksChan <- taskData:
		return nil // Envío exitoso
	default:
		return errors.New(MSG_TASK_SCHEDULER_BUSY_OR_CLOSED) // Canal lleno o cerrado, no bloquea
	}
}

// Agrega una tarea a ser ejecutada
func (scheduler *TaskScheduler) AddTask(task TaskFunc, mandatory bool, tag string) error {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	if scheduler.doHasTag(tag) {
		common.Log.Debugf(MSG_TAG_EXISTS, tag)
		return errors.New(MSG_TAG_EXISTS)
	}
	scheduler.taggedTasks[tag] = MAX_TASK_RETRIES
	common.Log.Debugf(MSG_TASK_ADDED, tag)
	err := scheduler.doAddTask(*newTaskDataWithoutExpirationTime(task, mandatory))
	return err
}

// Agrega una tarea a ser ejecutada con expiración en milisegundos
func (scheduler *TaskScheduler) AddTaskWithExpirationTime(task TaskFunc, deltaTime float32, tag string) error {
	scheduler.mutexTaggedTasks.Lock()
	defer scheduler.mutexTaggedTasks.Unlock()
	if scheduler.doHasTag(tag) {
		common.Log.Debugf(MSG_TAG_EXISTS, tag)
		return errors.New(MSG_TAG_EXISTS)
	}
	scheduler.taggedTasks[tag] = MAX_TASK_RETRIES
	common.Log.Debugf(MSG_TASK_ADDED, tag)
	err := scheduler.doAddTask(*newTaskDataWithExpirationTime(task, deltaTime))
	return err
}

// Remueve una etiqueta de la lista de tareas etiquetadas
func (scheduler *TaskScheduler) doRemoveTask(tag string) {
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
