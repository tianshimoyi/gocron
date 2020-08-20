package v1

import (
	"github.com/emicklei/go-restful"
	"github.com/x893675/gocron/internal/apiserver/models"
)

type systemHandler struct {
	hostModel models.HostStore
}

func newSystemHandler(hostModel models.HostStore) *systemHandler {
	return &systemHandler{hostModel: hostModel}
}

func (s *systemHandler) AddNode(request *restful.Request, response *restful.Response) {

}

func (s *systemHandler) DeleteNode(request *restful.Request, response *restful.Response) {

}

func (s *systemHandler) ListNode(request *restful.Request, response *restful.Response) {

}

func (s *systemHandler) GetNode(request *restful.Request, response *restful.Response) {

}

func (s *systemHandler) PingNode(request *restful.Request, response *restful.Response) {

}
