package v1

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/x893675/gocron/internal/apiserver/constants"
	"github.com/x893675/gocron/internal/apiserver/models"
	taskImpl "github.com/x893675/gocron/internal/apiserver/models/impl/task"
	taskLogImpl "github.com/x893675/gocron/internal/apiserver/models/impl/tasklog"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	taskSchema "github.com/x893675/gocron/internal/apiserver/schema"
	"github.com/x893675/gocron/internal/apiserver/service/task"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/server/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
)

const GroupName = "core"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

func AddToContainer(c *restful.Container, dbClient *database.Client, taskService *task.Task) error {
	ws := runtime.NewWebService(GroupVersion)
	handler := newTaskHandler(taskImpl.New(dbClient), taskLogImpl.New(dbClient), taskService)

	ws.Route(ws.POST("/tasks").
		To(handler.CreateTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Doc("Create task").
		Reads(taskSchema.TaskRequest{}).
		Returns(http.StatusCreated, constants.HTTP201, models.Task{}).
		Returns(http.StatusBadRequest, constants.HTTP400, restful.ServiceError{}).
		Returns(http.StatusConflict, constants.HTTP409, restful.ServiceError{}).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.GET("/tasks").
		To(handler.ListTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Doc("List task").
		Param(ws.QueryParameter(restplus.PagingParam, "paging query, e.g. limit=100,page=1").
			Required(false).
			DataFormat("limit=%d,page=%d").
			DefaultValue("limit=10,page=1")).
		Param(ws.QueryParameter(restplus.ReverseParam, "revers result").
			Required(false).
			DataType("bool").
			DefaultValue("false")).
		Param(ws.QueryParameter("id", "task id").
			Required(false).
			DataType("int").
			DefaultValue("0")).
		Param(ws.QueryParameter("hostid", "host id").
			Required(false).
			DataType("int").
			DefaultValue("0")).
		Param(ws.QueryParameter("tag", "task tag").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Param(ws.QueryParameter("name", "task name").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Param(ws.QueryParameter("protocol", "task protocol").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Param(ws.QueryParameter("status", "task status").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Param(ws.QueryParameter("creator", "task creator").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Writes([]models.Task{}).
		Returns(http.StatusOK, constants.HTTP200, restplus.PageableResponse{}).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.GET("/tasks/{task}").
		To(handler.GetTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("get task").
		Writes(models.Task{}).
		Returns(http.StatusOK, constants.HTTP200, models.Task{}).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.HEAD("/tasks/{task}").
		To(handler.CheckTaskExist).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id or name")).
		Doc("check task exist").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.DELETE("/tasks/{task}").
		To(handler.DeleteTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("delete task").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.PATCH("/tasks/{task}/enable").
		To(handler.EnableTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("enable task").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.PATCH("/tasks/{task}/disable").
		To(handler.DisableTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("disable task").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusInternalServerError, constants.HTTP500, restful.ServiceError{}))

	ws.Route(ws.POST("/tasks/{task}/run").
		To(handler.RunTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("run task").
		Returns(http.StatusOK, constants.HTTP200, nil))

	ws.Route(ws.POST("/tasks/{task}/stop").
		To(handler.StopTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("stop running task").
		Returns(http.StatusOK, constants.HTTP200, nil))

	ws.Route(ws.GET("/tasks/{task}/logs").
		To(handler.GetTaskLog).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("get task log").
		Param(ws.QueryParameter(restplus.PagingParam, "paging query, e.g. limit=100,page=1").
			Required(false).
			DataFormat("limit=%d,page=%d").
			DefaultValue("limit=10,page=1")).
		Param(ws.QueryParameter(restplus.ReverseParam, "revers result").
			Required(false).
			DataType("bool").
			DefaultValue("false")).
		Param(ws.QueryParameter("protocol", "task protocol").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Param(ws.QueryParameter("status", "task status").
			Required(false).
			DataType("string").
			DefaultValue("")).
		Returns(http.StatusOK, constants.HTTP200, []models.TaskLog{}))

	ws.Route(ws.DELETE("/tasks/{task}/logs").
		To(handler.ClearTaskLog).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("clear task log").
		Returns(http.StatusOK, constants.HTTP200, []models.TaskLog{}))

	c.Add(ws)
	return nil
}
