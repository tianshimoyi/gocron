package v1

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/x893675/gocron/internal/apiserver/constants"
	"github.com/x893675/gocron/internal/apiserver/models"
	taskImpl "github.com/x893675/gocron/internal/apiserver/models/impl/task"
	taskLogImpl "github.com/x893675/gocron/internal/apiserver/models/impl/tasklog"
	taskSchema "github.com/x893675/gocron/internal/apiserver/schema"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/server/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
)

const GroupName = "core"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

func AddToContainer(c *restful.Container, dbClient *database.Client) error {
	ws := runtime.NewWebService(GroupVersion)
	handler := newTaskHandler(taskImpl.New(dbClient), taskLogImpl.New(dbClient))

	ws.Route(ws.POST("/tasks").
		To(handler.CreateTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Doc("Create task").
		Reads(taskSchema.TaskRequest{}).
		Returns(http.StatusCreated, constants.HTTP201, models.Task{}))

	ws.Route(ws.GET("/tasks").
		To(handler.ListTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Doc("List task").
		Writes([]models.Task{}).
		Returns(http.StatusOK, constants.HTTP200, []models.Task{}))

	ws.Route(ws.GET("/tasks/{task}").
		To(handler.GetTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("get task").
		Writes(models.Task{}).
		Returns(http.StatusOK, constants.HTTP200, models.Task{}))

	ws.Route(ws.HEAD("/tasks/{task}").
		To(handler.GetTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("check task exist").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil))

	ws.Route(ws.DELETE("/tasks/{task}").
		To(handler.DeleteTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("delete task").
		Returns(http.StatusOK, constants.HTTP200, nil))

	ws.Route(ws.PATCH("/tasks/{task}/enable").
		To(handler.EnableTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("enable task").
		Returns(http.StatusOK, constants.HTTP200, nil))

	ws.Route(ws.PATCH("/tasks/{task}/disable").
		To(handler.DisableTask).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.TaskResourceTag}).
		Param(ws.PathParameter("task", "task id")).
		Doc("disable task").
		Returns(http.StatusOK, constants.HTTP200, nil))

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
