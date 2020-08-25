package task

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"github.com/x893675/gocron/internal/apiserver/models"
	taskImpl "github.com/x893675/gocron/internal/apiserver/models/impl/task"
	taskLogImpl "github.com/x893675/gocron/internal/apiserver/models/impl/tasklog"
	"github.com/x893675/gocron/pkg/client/database"
	"k8s.io/klog/v2"
	"runtime"
	"strconv"
	"sync"
)

const defaultConcurrencyQueue = 500

type Task struct {
	// 定时任务调度管理器
	serviceCron *cron.Cron

	// 同一任务是否有实例处于运行中
	runInstance *Instance

	// 任务计数-正在运行的任务
	taskCount *TaskCount

	// 并发队列, 限制同时运行的任务数量
	concurrencyQueue *ConcurrencyQueue

	//models task
	taskModel models.TaskStore

	//models task log
	taskLogModel models.TaskLogStore

	httpHandler  *HTTPHandler
	shellHandler *RPCHandler
}

func NewTaskService(dbClient *database.Client) *Task {
	t := &Task{
		serviceCron:      cron.New(),
		runInstance:      &Instance{m: sync.Map{}},
		taskCount:        &TaskCount{sync.WaitGroup{}, make(chan struct{})},
		concurrencyQueue: &ConcurrencyQueue{queue: make(chan struct{}, defaultConcurrencyQueue)},
		taskModel:        taskImpl.New(dbClient),
		taskLogModel:     taskLogImpl.New(dbClient),
		httpHandler:      &HTTPHandler{},
		shellHandler:     &RPCHandler{},
	}
	return t
}

// 初始化任务, 从数据库取出所有任务, 添加到定时任务并运行
func (t *Task) Initialize() {
	t.serviceCron.Start()
	go t.taskCount.Wait()
	klog.V(1).Infof("begin init corn job")
	//TODO: task model add
}

// 批量添加任务
func (t *Task) BatchAdd(tasks []models.Task) {
	for _, item := range tasks {
		t.RemoveAndAdd(item)
	}
}

// 删除任务后添加
func (t *Task) RemoveAndAdd(tasks models.Task) {
	//t.Remove(tasks.Id)
	//t.Add(tasks)
}

func (t *Task) Remove(id int) {
	t.serviceCron.RemoveJob(strconv.Itoa(id))
}

func (t *Task) NextRuntime() {

}

func (t *Task) Run() {

}

func (t *Task) Stop() {

}

func (t *Task) WaitAndExit() {

}

func (t *Task) ParseCronJobSpec(spec string) error {
	return PanicToError(func() {
		cron.Parse(spec)
	})
}

// 添加任务
func (t *Task) Add(taskModel models.Task) {
	//if taskModel.Level == models.TaskLevelChild {
	//	klog.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", taskModel.Id)
	//	return
	//}
	//taskFunc := t.createJob(taskModel)
	//if taskFunc == nil {
	//	klog.Error("创建任务处理Job失败,不支持的任务协议#", taskModel.Protocol)
	//	return
	//}
	//
	//cronName := strconv.Itoa(taskModel.Id)
	//err := PanicToError(func() {
	//	t.serviceCron.AddFunc(taskModel.Spec, taskFunc, cronName)
	//})
	//if err != nil {
	//	klog.Error("添加任务到调度器失败#", err)
	//}
}

//func (t *Task)createJob(taskModel models.Task) cron.FuncJob {
//	handler := createHandler(taskModel)
//	if handler == nil {
//		return nil
//	}
//	taskFunc := func() {
//		t.taskCount.Add()
//		defer t.taskCount.Done()
//
//		taskLogId := t.beforeExecJob(taskModel)
//		if taskLogId <= 0 {
//			return
//		}
//
//		if taskModel.Multi == 0 {
//			t.runInstance.add(taskModel.Id)
//			defer t.runInstance.done(taskModel.Id)
//		}
//
//		t.concurrencyQueue.Add()
//		defer t.concurrencyQueue.Done()
//
//		klog.V(1).Infof("开始执行任务#%s#命令-%s", taskModel.Name, taskModel.Command)
//		taskResult := execJob(handler, taskModel, taskLogId)
//		klog.V(1).Infof("任务完成#%s#命令-%s", taskModel.Name, taskModel.Command)
//		afterExecJob(taskModel, taskResult, taskLogId)
//	}
//
//	return taskFunc
//}
//
//// 任务前置操作
//func (t *Task) beforeExecJob(taskModel models.Task) (taskLogId int64) {
//	if taskModel.Multi == 0 && t.runInstance.has(taskModel.Id) {
//		createTaskLog(taskModel, models.Cancel)
//		return
//	}
//	taskLogId, err := createTaskLog(taskModel, models.Running)
//	if err != nil {
//		klog.Error("任务开始执行#写入任务日志失败-", err)
//		return
//	}
//	klog.V(1).Infof("任务命令-%s", taskModel.Command)
//
//	return taskLogId
//}

func createHandler(taskModel models.Task) Handler {
	var handler Handler = nil
	switch taskModel.Protocol {
	case models.TaskHTTP:
		handler = new(HTTPHandler)
	case models.TaskShell:
		handler = new(RPCHandler)
	}

	return handler
}

type Handler interface {
	Run(taskModel models.Task, taskUniqueId int64) (string, error)
}

type TaskResult struct {
	Result     string
	Err        error
	RetryTimes int8
}

// PanicToError Panic转换为error
func PanicToError(f func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(PanicTrace(e))
		}
	}()
	f()
	return
}

// PanicTrace panic调用链跟踪
func PanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)

	return fmt.Sprintf("panic: %v %s", err, stackBuf[:n])
}
