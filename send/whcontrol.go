package send

import (
	"github.com/robertzml/Gorgons/tlv"
	"strconv"
)

// 设备控制报文
// 0x10
type whControlMessage struct {
	SerialNumber    string
	MainboardNumber string
	ControlAction   string
}

// 设备控制报文构造函数
func NewWHControlMessage(serialNumber string, mainboardNumber string) *whControlMessage {
	return &whControlMessage{SerialNumber: serialNumber, MainboardNumber: mainboardNumber}
}

// 拼接设备控制报文
func (msg *whControlMessage) splice() string {
	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	ca := tlv.Splice(0x012, msg.ControlAction)

	body := tlv.Splice(0x0010, sn+mn+ca)

	return head + body
}

// 开关机报文
func (msg *whControlMessage) Power(power int) string {
	msg.ControlAction = tlv.Splice(0x01, strconv.Itoa(power))
	return msg.splice()
}

// 激活非激活报文
func (msg *whControlMessage) Activate(status int) string {
	msg.ControlAction = tlv.Splice(0x1b, strconv.Itoa(status))
	return msg.splice()
}

// 设备加锁报文
func (msg *whControlMessage) Lock() string {
	msg.ControlAction = tlv.Splice(0x1a, strconv.Itoa(0))
	return msg.splice()
}

// 设备解锁报文
func (msg *whControlMessage) Unlock(option int, deadline int64) string {
	unlock := tlv.Splice(0x1a, strconv.Itoa(1))

	if option == 1 {
		dl := tlv.ParseTimestampToString(deadline)
		msg.ControlAction = unlock + tlv.Splice(0x20, dl)
	} else {
		msg.ControlAction = unlock
	}

	return msg.splice()
}
