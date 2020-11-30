package main

import (
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/glog"
	"github.com/robertzml/Gorgons/mqtt"
	"github.com/robertzml/Gorgons/queue"
	"github.com/robertzml/Gorgons/redis"
)

func main() {
	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()

	// 载入配置文件
	base.LoadConfig()

	// 初始化全局channel
	base.InitChannel()

	// 启动日志服务
	glog.InitGlog()
	go startLog()

	// 初始化redis连接池
	redisClient := redis.Init()

	// 启动 MQTT订阅服务
	mqtt.InitMQTT()
	go startControl()

	// 初始化队列服务
	_ = queue.InitQueue(redisClient)

	// 启动接收数据处理
	startPipe()

	// 阻塞
	select{}
}

// 启动日志服务
func startLog() {
	fmt.Println("start log service.")
	glog.Read()
}

// 启动设备控制服务
func startControl() {
	fmt.Println("start control service.")
	mqtt.StartSend()
}

// 启动接收数据处理
func startPipe() {
	fmt.Println("start queue service.")
	go queue.Control()
	go queue.Feedback()
	go queue.Special()
}