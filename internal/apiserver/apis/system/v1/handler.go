package v1

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	"github.com/x893675/gocron/internal/apiserver/rpc"
	"github.com/x893675/gocron/internal/apiserver/schema"
	"github.com/x893675/gocron/pkg/pb"
	"github.com/x893675/gocron/pkg/utils/stringutils"
	"net/http"
)

const testConnectionCommand = "echo hello"
const testConnectionTimeout = 5

type systemHandler struct {
	hostModel models.HostStore
}

func newSystemHandler(hostModel models.HostStore) *systemHandler {
	return &systemHandler{hostModel: hostModel}
}

func (s *systemHandler) AddNode(request *restful.Request, response *restful.Response) {
	var host schema.NodeRequest
	err := restplus.ParseBody(request, &host)
	if err != nil {
		restplus.HandleBadRequest(response, request, err)
		return
	}
	exist, err := s.hostModel.Exist(request.Request.Context(), models.GetParam{
		Name: host.Name,
	})
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	if exist {
		restplus.HandleConflict(response, request, fmt.Errorf("conflict node name"))
		return
	}
	err = s.hostModel.Create(request.Request.Context(), models.SchemaHost(host).ToModelHost())
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	response.WriteHeader(http.StatusCreated)
}

func (s *systemHandler) DeleteNode(request *restful.Request, response *restful.Response) {
	node := request.PathParameter("node")
	var param models.GetParam
	parseIdOrName(node, &param)
	err := s.hostModel.Delete(request.Request.Context(), models.DeleteParam(param))
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func (s *systemHandler) ListNode(request *restful.Request, response *restful.Response) {
	limit, offset := restplus.ParsePaging(request)
	reverse := restplus.GetBoolValueWithDefault(request, restplus.ReverseParam, false)
	hosts, total, err := s.hostModel.List(request.Request.Context(), models.ListHostParam{
		BaseListParam: models.BaseListParam{
			Limit:   limit,
			Offset:  offset,
			Reverse: reverse,
		},
	})
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	restplus.ResWithPage(response, hosts, int(total), http.StatusOK)
}

func (s *systemHandler) GetNode(request *restful.Request, response *restful.Response) {
	node := request.PathParameter("node")
	var param models.GetParam
	parseIdOrName(node, &param)
	host, err := s.hostModel.Get(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	_ = response.WriteHeaderAndEntity(http.StatusOK, host)
}

func (s *systemHandler) CheckNodeExist(request *restful.Request, response *restful.Response) {
	node := request.PathParameter("node")
	var param models.GetParam
	parseIdOrName(node, &param)
	exist, err := s.hostModel.Exist(request.Request.Context(), param)
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

func (s *systemHandler) UpdateNode(request *restful.Request, response *restful.Response) {
	node := request.PathParameter("node")
	var param models.GetParam
	parseIdOrName(node, &param)
	oldHost, err := s.hostModel.Get(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	if oldHost.Name == "" {
		restplus.HandleNotFound(response, request, fmt.Errorf("node %s not exist", node))
		return
	}
	var newHost schema.NodeRequest
	err = restplus.ParseBody(request, &newHost)
	if err != nil {
		restplus.HandleBadRequest(response, request, err)
		return
	}
	if newHost.Name != oldHost.Name {
		exist, err := s.hostModel.Exist(request.Request.Context(), models.GetParam{
			Name: newHost.Name,
		})
		if err != nil {
			restplus.HandleInternalError(response, request, err)
			return
		}
		if exist {
			restplus.HandleConflict(response, request, fmt.Errorf("conflict node name"))
			return
		}
	}
	item := models.SchemaHost(newHost).ToModelHost()
	item.Id = oldHost.Id
	err = s.hostModel.Update(request.Request.Context(), item)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	response.WriteHeader(http.StatusOK)
}

func (s *systemHandler) PingNode(request *restful.Request, response *restful.Response) {
	node := request.PathParameter("node")
	var param models.GetParam
	parseIdOrName(node, &param)
	host, err := s.hostModel.Get(request.Request.Context(), param)
	if err != nil {
		restplus.HandleInternalError(response, request, err)
		return
	}
	if host.Name == "" {
		restplus.HandleBadRequest(response, request, fmt.Errorf("node not exsit"))
		return
	}
	output, err := rpc.Exec(host.Addr, host.Port, &pb.TaskRequest{
		Command: testConnectionCommand,
		Timeout: testConnectionTimeout,
	})
	if err != nil {
		restplus.HandleExpectedFailed(response, request, fmt.Errorf("connect faild: %v %v", err, output))
		return
	}
	response.WriteHeader(http.StatusOK)
}

func parseIdOrName(param string, out *models.GetParam) {
	id := stringutils.S(param).DefaultInt(0)
	if id > 0 {
		out.ID = id
	} else {
		out.Name = param
	}
}
