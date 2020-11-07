package main

import (
	"fmt"
	"github.com/robertzml/Gorgons/base"
	"github.com/robertzml/Gorgons/pipe"
)

func main() {
	fmt.Println("app is running")

	defer func() {
		fmt.Println("app is stop.")
	}()

	// 载入配置文件
	base.LoadConfig()

	// 启动接收数据处理
	go startPipe()

	// 阻塞
	select{}
}

// 启动接收数据处理
func startPipe() {
	fmt.Println("start pipe service")
	pipe.Process()
}