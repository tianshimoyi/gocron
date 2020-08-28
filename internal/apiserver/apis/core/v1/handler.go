package v1

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	"github.com/x893675/gocron/internal/apiserver/schema"
	"github.com/x893675/gocron/internal/apiserver/service/task"
	"github.com/x893675/gocron/pkg/utils/stringutils"
	"k8s.io/klog/v2"
	"net/http"
)

type taskHandler struct {
	taskModel    models.TaskStore
	taskLogModel models.TaskLogStore
	taskService  *task.Task
}

func newTaskHandler(taskModel models.TaskStore, taskLogModel models.TaskLogStore, taskService *task.Task) *taskHandler {
	return &taskHandler{taskModel: taskModel, taskLogModel: taskLogModel, taskService: taskService}
}

func (t *taskHandler) CreateTask(request *restful.Request, response *restful.Response) {
	var item schema.TaskRequest
	err := restplus.ParseBody(request, &item)
	if err != nil {
		restplus.HandleBadRequest(response, request, err)
		return
	}
	exist, err := t.taskModel.Exist(request.Request.Context(), models.GetParam{
		Name: item.Name,
	})
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	if exist {
		restplus.HandleConflict(response, request, fmt.Errorf("conflict task name"))
		return
	}
	if item.Protocol == models.TaskShell && item.HostId == "" {
		restplus.HandleBadRequest(response, request, fmt.Errorf("host id must be valid when task protocol is shell"))
		return
	}
	if item.Type == models.TaskTypeCronJob {
		if err := t.taskService.ParseCronJobSpec(item.Spec); err != nil {
			restplus.HandleBadRequest(response, request, err)
			return
		}
	}
	taskID, err := t.taskModel.Create(request.Request.Context(), models.SchemaTask(item))
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	//TODO: add run at job
	switch item.Type {
	case models.TaskTypeCronJob:
		t.addTaskToTimer(int(taskID))
	case models.TaskTypeJob:
		fallthrough
	default:
		t.runJob(int(taskID))
	}
	response.WriteHeader(http.StatusCreated)
}

func (t *taskHandler) GetTask(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	var param models.GetParam
	parseIdOrName(tsk, &param)
	result, err := t.taskModel.Get(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	_ = response.WriteHeaderAndEntity(http.StatusOK, result)
}

func (t *taskHandler) ListTask(request *restful.Request, response *restful.Response) {
	limit, offset := restplus.ParsePaging(request)

	param := models.ListTaskParam{
		BaseListParam: models.BaseListParam{
			Reverse: restplus.GetBoolValueWithDefault(request, restplus.ReverseParam, false),
			Offset:  offset,
			Limit:   limit,
		},
		GetParam: models.GetParam{
			ID: restplus.GetIntValueWithDefault(request, "id", 0),
		},
		Status:   restplus.GetStringValueWithDefault(request, "status", ""),
		Level:    restplus.GetStringValueWithDefault(request, "level", ""),
		HostID:   restplus.GetIntValueWithDefault(request, "hostid", 0),
		Protocol: restplus.GetStringValueWithDefault(request, "protocol", ""),
		Tag:      restplus.GetStringValueWithDefault(request, "tag", ""),
	}
	result, total, err := t.taskModel.List(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	for _, job := range result {
		t.taskService.NextRuntime(job)
	}
	restplus.ResWithPage(response, result, int(total), http.StatusOK)
}

func (t *taskHandler) DeleteTask(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	var param models.GetParam
	parseIdOrName(tsk, &param)
	if param.ID == 0 {
		restplus.HandleBadRequest(response, request, fmt.Errorf("only support delete task by id"))
		return
	}
	err := t.taskModel.Delete(request.Request.Context(), models.DeleteParam(param))
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	t.taskService.Remove(param.ID)
	response.WriteHeader(http.StatusOK)
}

func (t *taskHandler) EnableTask(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	id := stringutils.S(tsk).DefaultInt(0)
	if id == 0 {
		restplus.HandleBadRequest(response, request, fmt.Errorf("only support delete task by id"))
		return
	}
	err := t.taskModel.UpdateTaskStatus(request.Request.Context(), id, models.TaskStatusEnabled)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	t.addTaskToTimer(id)
	response.WriteHeader(http.StatusOK)
}

func (t *taskHandler) DisableTask(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	id := stringutils.S(tsk).DefaultInt(0)
	if id == 0 {
		restplus.HandleBadRequest(response, request, fmt.Errorf("only support delete task by id"))
		return
	}
	err := t.taskModel.UpdateTaskStatus(request.Request.Context(), id, models.TaskStatusDisabled)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	t.taskService.Remove(id)
	response.WriteHeader(http.StatusOK)
}

func (t *taskHandler) RunTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) GetTaskLog(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	id := stringutils.S(tsk).DefaultInt(0)
	if id == 0 {
		restplus.HandleBadRequest(response, request, fmt.Errorf("only support get task log by id"))
		return
	}
	limit, offset := restplus.ParsePaging(request)
	result, total, err := t.taskLogModel.List(request.Request.Context(), models.ListTaskLogParam{
		BaseListParam: models.BaseListParam{
			Limit:   limit,
			Offset:  offset,
			Reverse: restplus.GetBoolValueWithDefault(request, restplus.ReverseParam, false),
		},
		Status:   restplus.GetStringValueWithDefault(request, "status", ""),
		Protocol: restplus.GetStringValueWithDefault(request, "protocol", ""),
		TaskID:   id,
	})
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	restplus.ResWithPage(response, result, int(total), http.StatusOK)
}

func (t *taskHandler) ClearTaskLog(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	id := stringutils.S(tsk).DefaultInt(0)
	if id == 0 {
		restplus.HandleBadRequest(response, request, fmt.Errorf("only support delete task log by id"))
		return
	}
	if err := t.taskLogModel.Delete(request.Request.Context(), models.DeleteTaskLogParam{
		TaskID: id,
		Status: restplus.GetStringValueWithDefault(request, "status", ""),
		Mon:    restplus.GetIntValueWithDefault(request, "mon", 0),
	}); err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func (t *taskHandler) StopTask(request *restful.Request, response *restful.Response) {

}

func (t *taskHandler) CheckTaskExist(request *restful.Request, response *restful.Response) {
	tsk := request.PathParameter("task")
	var param models.GetParam
	parseIdOrName(tsk, &param)
	exist, err := t.taskModel.Exist(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	if exist {
		response.WriteHeader(http.StatusOK)
		return
	}
	response.WriteHeader(http.StatusNotFound)
}

func parseIdOrName(param string, out *models.GetParam) {
	id := stringutils.S(param).DefaultInt(0)
	if id > 0 {
		out.ID = id
	} else {
		out.Name = param
	}
}

func (t *taskHandler) addTaskToTimer(id int) {
	job, err := t.taskModel.Get(context.TODO(), models.GetParam{
		ID: id,
	})
	if err != nil {
		klog.Error(err)
		return
	}
	t.taskService.RemoveAndAdd(job)
}

func (t *taskHandler) runJob(id int) {
	job, err := t.taskModel.Get(context.TODO(), models.GetParam{
		ID: id,
	})
	if err != nil {
		klog.Error(err)
		return
	}
	t.taskService.Run(job)
}
