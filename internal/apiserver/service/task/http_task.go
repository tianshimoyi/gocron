package task

import (
	"fmt"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/pkg/client/httpclient"
	"net/http"
	"strings"
)

// HTTP任务
type HTTPHandler struct{}

// http任务执行时间不超过300秒
const HttpExecTimeout = 300

func (h *HTTPHandler) Run(taskModel *models.Task, taskUniqueId int64) (result string, err error) {
	if taskModel.Timeout <= 0 || taskModel.Timeout > HttpExecTimeout {
		taskModel.Timeout = HttpExecTimeout
	}
	var resp httpclient.ResponseWrapper
	if taskModel.HttpMethod == models.TaskHTTPMethodGet {
		resp = httpclient.Get(taskModel.Command, taskModel.Timeout)
	} else {
		urlFields := strings.Split(taskModel.Command, "?")
		taskModel.Command = urlFields[0]
		var params string
		if len(urlFields) >= 2 {
			params = urlFields[1]
		}
		resp = httpclient.PostParams(taskModel.Command, params, taskModel.Timeout)
	}
	// 返回状态码非200，均为失败
	if resp.StatusCode != http.StatusOK {
		return resp.Body, fmt.Errorf("HTTP状态码非200-->%d", resp.StatusCode)
	}

	return resp.Body, err
}
