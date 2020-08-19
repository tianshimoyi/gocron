package v1

import (
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/models"
)

type taskHandler struct {
	taskModel    models.TaskStore
	taskLogModel models.TaskLogStore
}

func newTaskHandler(taskModel models.TaskStore, taskLogModel models.TaskLogStore) *taskHandler {
	return &taskHandler{taskModel: taskModel, taskLogModel: taskLogModel}
}

func (t *taskHandler) CreateTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) GetTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) ListTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) DeleteTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) EnableTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) DisableTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) RunTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) GetTaskLog(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) ClearTaskLog(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) StopTask(request *restful.Request, response *restful.Response) {

}
