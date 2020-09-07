# gocron

使用Go语言开发的轻量级定时任务集中调度和管理系统

## 功能特性

* 支持 `cronjob` (crontab时间表达式), `job` ( run once, 创建成功立即执行，且只执行一次 ), `planjob` ( 计划任务，在指定时间执行任务，且只执行一次 )
* 任务执行失败可重试
* 任务执行超时，强制结束
* 支持 `shell` 和 `http` 任务
* 支持在多节点执行任务
* 支持任务日志

## 持久化数据库

* mysql
* postgres

## 数据表结构

![data schema](docs/imgs/gocron-schema.png)

## 部署

### docker-compose部署

1. `docker-compose up -d`

### k8s部署

1. `kubectl create ns gocron`
2. `kubectl apply -f deployment/deploy.yaml -n gocron`

## Api测试例子

查看代码跟目录 `gocron.postman_collection.json`文件，导入 `postman` 等工具