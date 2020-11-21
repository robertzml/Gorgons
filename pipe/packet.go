package pipe

// 队列发送包类型
type queuePacket struct {
	SerialNumber	string
	DeviceType		int
	ControlType		int
	Option 			int
	Deadline		int64
}