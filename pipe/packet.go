package pipe

// 队列发送包类型
type packet struct {
	SerialNumber	string
	DeviceType		int
	ControlType		int
	Option 			int
	Deadline		int64
}