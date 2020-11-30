package queue

// 队列设备控制包类型
type controlPacket struct {
	// 序列号
	SerialNumber	string

	// 设备类型
	DeviceType		int

	// 控制类型
	ControlType		int

	// 控制参数
	Option 			int

	// 允许使用时间
	Deadline		int64
}

// 队列设备反馈包类型
type feedbackPacket struct {
	// 序列号
	SerialNumber	string

	// 设备类型
	DeviceType		int

	// 控制类型
	ControlType		int

	// 控制参数
	Option 			int
}

// 队列设备特殊包类型
type specialPacket struct {
	// 序列号
	SerialNumber	string

	// 设备类型
	DeviceType		int

	// 控制类型
	ControlType		int

	// 控制参数
	Option 			string
}