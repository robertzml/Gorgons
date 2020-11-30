package queue

import (
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/db"
	"github.com/robertzml/Gorgons/glog"
	"github.com/streadway/amqp"
)

const (
	// 当前包名称
	packageName = "queue"
)

var (
	// 用于注入实时数据访问类
	snapshot db.Snapshot

	// rabbit mq 连接
	rmConnection *amqp.Connection
)

/**
初始化队列服务
 */
func InitQueue(snap db.Snapshot) (err error){
	snapshot = snap

	rmConnection, err = amqp.Dial(base.DefaultConfig.RabbitMQAddress)
	if err != nil {
		glog.Write(1, packageName, "InitQueue", "connect to rabbit mq failed.")
	}

	return err
}