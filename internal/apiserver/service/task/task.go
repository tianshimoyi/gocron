package task

import (
	"context"
	"fmt"
	"github.com/jakecoffman/cron"
	"github.com/x893675/gocron/internal/apiserver/models"
	taskImpl "github.com/x893675/gocron/internal/apiserver/models/impl/task"
	taskLogImpl "github.com/x893675/gocron/internal/apiserver/models/impl/tasklog"
	"github.com/x893675/gocron/internal/apiserver/rpc"
	"github.com/x893675/gocron/pkg/client/database"
	"k8s.io/klog/v2"
	"runtime"
	"sync"
	"time"
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
func (t *Task) Initialize(stopCh <-chan struct{}) error {
	t.serviceCron.Start()
	go func() {
		<-stopCh
		t.WaitAndExit()
	}()
	go t.taskCount.Wait()
	klog.V(1).Infof("begin init corn job")
	maxPage := 1000
	pageSize := 1000
	page := 1
	taskNum := 0
	for page < maxPage {
		taskList, total, err := t.taskModel.List(context.TODO(), models.ListTaskParam{
			Status: models.TaskStatusEnabled,
			Level:  models.ParentLevelTask,
			Type:   models.TaskTypeCronJob,
			BaseListParam: models.BaseListParam{
				Offset: (page - 1) * pageSize,
				Limit:  pageSize,
			},
		})
		if err != nil {
			return err
		}
		for _, item := range taskList {
			t.Add(item)
			taskNum++
		}
		maxPage = int(total) % pageSize
		page++
	}
	klog.V(1).Infof("%d corn job init finish", taskNum)
	return nil
}

// 批量添加任务
func (t *Task) BatchAdd(tasks []models.Task) {
	//for _, item := range tasks {
	//	t.RemoveAndAdd(item)
	//}
}

// 删除任务后添加
func (t *Task) RemoveAndAdd(job *models.Task) {
	t.Remove(int(job.Id))
	t.Add(job)
}

func (t *Task) Remove(id int) {
	t.serviceCron.RemoveJob(fmt.Sprintf("task-%d", id))
}

func (t *Task) NextRuntime(job *models.Task) {
	if job.Status != models.TaskStatusEnabled || job.Level != models.ParentLevelTask {
		return
	}
	e := t.serviceCron.Entries()
	for _, item := range e {
		if item.Name == fmt.Sprintf("task-%d", job.Id) {
			ts := item.Next
			job.NextRunTime = &ts
			break
		}
	}
}

func (t *Task) Run(job *models.Task) {
	go t.createJob(job)()
}

func (t *Task) Stop(ip string, port int, id int64) {
	rpc.Stop(ip, port, id)
}

func (t *Task) WaitAndExit() {
	t.serviceCron.Stop()
	t.taskCount.Exit()
}

func (t *Task) ParseCronJobSpec(spec string) error {
	return PanicToError(func() {
		cron.Parse(spec)
	})
}

// 添加任务
func (t *Task) Add(job *models.Task) {
	if job.Level == models.ChildLevelTask {
		klog.Errorf("添加任务失败#不允许添加子任务到调度器#任务Id-%d", job.Id)
		return
	}
	taskFunc := t.createJob(job)
	if taskFunc == nil {
		klog.Error("创建任务处理Job失败,不支持的任务协议#", job.Protocol)
		return
	}
	cronName := fmt.Sprintf("task-%d", job.Id)
	err := PanicToError(func() {
		t.serviceCron.AddFunc(job.Spec, taskFunc, cronName)
	})
	if err != nil {
		klog.Error("添加任务到调度器失败#", err)
	}
}

func (t *Task) createJob(job *models.Task) cron.FuncJob {
	switch job.Protocol {
	case models.TaskShell:
		break
	case models.TaskHTTP:
		break
	default:
		return nil
	}
	taskFunc := func() {
		t.taskCount.Add()
		defer t.taskCount.Done()
		taskLogId := t.beforeExecJob(job)
		if taskLogId <= 0 {
			return
		}
		if !job.Multi {
			t.runInstance.add(int(job.Id))
			defer t.runInstance.done(int(job.Id))
		}
		t.concurrencyQueue.Add()
		defer t.concurrencyQueue.Done()
		klog.V(2).Infof("开始执行任务#%s#命令-%s", job.Name, job.Command)
		taskResult := t.execJob(job, taskLogId)
		klog.V(2).Infof("任务完成#%s#命令-%s", job.Name, job.Command)
		t.afterExecJob(job, taskResult, taskLogId)
	}
	return taskFunc
}

func (t *Task) execJob(job *models.Task, logID int64) TaskResult {
	defer func() {
		if err := recover(); err != nil {
			klog.Error("panic#service/task.go:execJob#", err)
		}
	}()
	// 默认只运行任务一次
	var execTimes int8 = 1
	if job.RetryTimes > 0 {
		execTimes += job.RetryTimes
	}
	var i int8 = 0
	var output string
	var err error
	for i < execTimes {
		switch job.Protocol {
		case models.TaskHTTP:
			output, err = t.httpHandler.Run(job, logID)
		case models.TaskShell:
			fallthrough
		default:
			output, err = t.shellHandler.Run(job, logID)
		}
		if err == nil {
			return TaskResult{Result: output, Err: err, RetryTimes: i}
		}
		i++
		if i < execTimes {
			klog.Warningf("任务执行失败#任务id-%d#重试第%d次#输出-%s#错误-%s", job.Id, i, output, err.Error())
			if job.RetryInterval > 0 {
				time.Sleep(time.Duration(job.RetryInterval) * time.Second)
			} else {
				// 默认重试间隔时间，每次递增1分钟
				time.Sleep(time.Duration(i) * time.Minute)
			}
		}
	}
	return TaskResult{Result: output, Err: err, RetryTimes: job.RetryTimes}
}

func (t *Task) beforeExecJob(job *models.Task) (taskLogId int64) {
	if !job.Multi && t.runInstance.has(int(job.Id)) {
		//FIXME
		_, _ = t.createTaskLog(job, models.TaskLogStatusCancel)
		return
	}
	taskLogId, err := t.createTaskLog(job, models.TaskLogStatusRunning)
	if err != nil {
		klog.Error("任务开始执行#写入任务日志失败-", err)
		return
	}
	klog.V(2).Infof("任务命令-%s", job.Command)
	return taskLogId
}

func (t *Task) afterExecJob(job *models.Task, taskResult TaskResult, taskLogId int64) {
	err := t.updateTaskLog(taskLogId, taskResult)
	if err != nil {
		klog.Error("任务结束#更新任务日志失败-", err)
	}

	//TODO
	// 发送邮件
	//go SendNotification(taskModel, taskResult)
	// 执行依赖任务
	//go execDependencyTask(taskModel, taskResult)
}

// 更新任务日志
func (t *Task) updateTaskLog(taskLogId int64, taskResult TaskResult) error {
	taskLog := &models.TaskLog{
		Id:         uint(taskLogId),
		RetryTimes: taskResult.RetryTimes,
		Result:     taskResult.Result,
	}
	if taskResult.Err != nil {
		taskLog.Status = models.TaskLogStatusFailure
	} else {
		taskLog.Status = models.TaskLogStatusFinish
	}
	return t.taskLogModel.Update(context.TODO(), taskLog)
}

func (t *Task) createTaskLog(job *models.Task, status string) (int64, error) {
	taskLog := &models.TaskLog{
		TaskId:     int(job.Id),
		Name:       job.Name,
		Spec:       job.Spec,
		Protocol:   job.Protocol,
		Command:    job.Command,
		Timeout:    job.Timeout,
		RetryTimes: 0,
		Status:     status,
	}
	if job.Protocol == models.TaskShell {
		aggregationHost := ""
		for _, host := range job.Hosts {
			aggregationHost += fmt.Sprintf("%s - %s<br>", host.Alias, host.Name)
		}
		taskLog.Hostname = aggregationHost
	}
	return t.taskLogModel.Create(context.TODO(), taskLog)
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

//func createHandler(taskModel models.Task) Handler {
//	var handler Handler = nil
//	switch taskModel.Protocol {
//	case models.TaskHTTP:
//		handler = new(HTTPHandler)
//	case models.TaskShell:
//		handler = new(RPCHandler)
//	}
//
//	return handler
//}

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
