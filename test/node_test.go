package test

import (
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/internal/apiserver/restplus"
	testUtil "github.com/x893675/gocron/test/utils"
	"net/http"
	"testing"
)

var (
	addNodeURL  = "/api/system/v1/nodes"
	addNodeBody = `{
    "name": "node1",
    "alias": "test",
    "port": 8090,
    "addr": "localhost"
}`
	getNodeURL        = "/api/system/v1/nodes/node1"
	ListNodeURL       = addNodeURL
	checkNodeExistURL = getNodeURL
	deleteNodeURL     = "/api/system/v1/nodes/1"
)

func testAddNode(t *testing.T) {
	_, httpStatus := requestTest(t, getRequestURL(addNodeURL), http.MethodPost, []byte(addNodeBody), nil, true, true)
	testUtil.HttpStatusIsExpected(t, []int{http.StatusCreated}, httpStatus)
}

func testGetNode(t *testing.T) {
	body, httpStatus := requestTest(t, getRequestURL(getNodeURL), http.MethodGet, nil, nil, true, true)
	ret := models.Host{}
	if !testUtil.TryJsonUnMarsh(t, body, &ret, "models.Host") {
		return
	}
	testUtil.HttpStatusIsExpected(t, []int{http.StatusOK}, httpStatus)
}

func testListNode(t *testing.T) {
	body, httpStatus := requestTest(t, getRequestURL(ListNodeURL), http.MethodGet, nil, nil, true, true)
	ret := restplus.PageableResponse{}
	if !testUtil.TryJsonUnMarsh(t, body, &ret, "restplus.PageableResponse") {
		return
	}
	testUtil.HttpStatusIsExpected(t, []int{http.StatusOK}, httpStatus)
	if ret.TotalCount != 1 {
		t.Errorf("expected node number is 1, bug got %v", ret.TotalCount)
	}
}

func testCheckNodeExist(t *testing.T) {
	_, httpStatus := requestTest(t, getRequestURL(checkNodeExistURL), http.MethodHead, nil, nil, true, true)
	testUtil.HttpStatusIsExpected(t, []int{http.StatusOK}, httpStatus)
}

func testCheckNodeNotExist(t *testing.T) {
	_, httpStatus := requestTest(t, getRequestURL(checkNodeExistURL), http.MethodHead, nil, nil, true, true)
	testUtil.HttpStatusIsExpected(t, []int{http.StatusNotFound}, httpStatus)
}

func testDeleteNode(t *testing.T) {
	_, httpStatus := requestTest(t, getRequestURL(deleteNodeURL), http.MethodDelete, nil, nil, true, true)
	testUtil.HttpStatusIsExpected(t, []int{http.StatusOK}, httpStatus)
}

func TestNodeApi(t *testing.T) {
	t.Run("Check node1 not exsit", testCheckNodeNotExist)
	t.Run("Create node1", testAddNode)
	t.Run("Check node1 exsit", testCheckNodeExist)
	t.Run("List nodes", testListNode)
	t.Run("Get node1", testGetNode)
	t.Run("DeleteNode", testDeleteNode)
}
