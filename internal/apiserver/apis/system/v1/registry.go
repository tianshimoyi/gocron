package v1

import (
	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/x893675/gocron/internal/apiserver/constants"
	"github.com/x893675/gocron/internal/apiserver/models"
	hostImpl "github.com/x893675/gocron/internal/apiserver/models/impl/node"
	taskSchema "github.com/x893675/gocron/internal/apiserver/schema"
	"github.com/x893675/gocron/pkg/client/database"
	"github.com/x893675/gocron/pkg/server/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"net/http"
)

const GroupName = "system"

var GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1"}

func AddToContainer(c *restful.Container, dbClient *database.Client) error {

	ws := runtime.NewWebService(GroupVersion)
	handler := newSystemHandler(hostImpl.New(dbClient))

	ws.Route(ws.POST("/nodes").
		To(handler.AddNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Doc("Add host").
		Reads(taskSchema.TaskRequest{}).
		Returns(http.StatusCreated, constants.HTTP201, models.Host{}))

	ws.Route(ws.GET("/nodes").
		To(handler.ListNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Doc("List host").
		Writes([]models.Host{}).
		Returns(http.StatusOK, constants.HTTP200, []models.Host{}))

	ws.Route(ws.GET("/nodes/{node}").
		To(handler.GetNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Param(ws.PathParameter("node", "node id")).
		Doc("Get host").
		Writes(models.Host{}).
		Returns(http.StatusOK, constants.HTTP200, models.Host{}))

	ws.Route(ws.HEAD("/nodes/{node}").
		To(handler.GetNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Doc("Check host exist").
		Param(ws.PathParameter("node", "node id")).
		Writes(models.Host{}).
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil))

	ws.Route(ws.DELETE("/nodes/{node}").
		To(handler.DeleteNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Param(ws.PathParameter("node", "node id")).
		Doc("Delete host").
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil))

	ws.Route(ws.GET("/nodes/{node}/ping").
		To(handler.PingNode).
		Metadata(restfulspec.KeyOpenAPITags, []string{constants.NodeResourceTag}).
		Doc("Check host exist").
		Param(ws.PathParameter("node", "node id")).
		Writes(models.Host{}).
		Returns(http.StatusOK, constants.HTTP200, nil).
		Returns(http.StatusNotFound, constants.HTTP404, nil))

	c.Add(ws)
	return nil
}
