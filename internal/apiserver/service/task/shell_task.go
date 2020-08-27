package task

import (
	"fmt"
	"github.com/x893675/gocron/internal/apiserver/models"
	"github.com/x893675/gocron/internal/apiserver/rpc"
	"github.com/x893675/gocron/pkg/pb"
)

// RPC调用执行任务
type RPCHandler struct{}

func (h *RPCHandler) Run(taskModel *models.Task, taskUniqueId int64) (result string, err error) {
	taskRequest := &pb.TaskRequest{}
	taskRequest.Timeout = int32(taskModel.Timeout)
	taskRequest.Command = taskModel.Command
	taskRequest.Id = taskUniqueId
	resultChan := make(chan TaskResult, len(taskModel.Hosts))
	for _, taskHost := range taskModel.Hosts {
		go func(th models.TaskHostDetail) {
			output, err := rpc.Exec(th.Addr, th.Port, taskRequest)
			errorMessage := ""
			if err != nil {
				errorMessage = err.Error()
			}
			outputMessage := fmt.Sprintf("主机: [%s-%s:%d]\n%s\n%s\n\n",
				th.Alias, th.Name, th.Port, errorMessage, output,
			)
			resultChan <- TaskResult{Err: err, Result: outputMessage}
		}(taskHost)
	}

	var aggregationErr error = nil
	aggregationResult := ""
	for i := 0; i < len(taskModel.Hosts); i++ {
		taskResult := <-resultChan
		aggregationResult += taskResult.Result
		if taskResult.Err != nil {
			aggregationErr = taskResult.Err
		}
	}

	return aggregationResult, aggregationErr
}
