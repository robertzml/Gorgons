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

// 设定温度报文
func (msg *whControlMessage) SetTemp(temp int) string {
	msg.ControlAction = tlv.Splice(0x1c, strconv.FormatInt(int64(temp), 16))
	return msg.splice()
}

// 设置允许使用时间
func (msg *whControlMessage) SetDeadline(deadline int64) string {
	dl := tlv.ParseTimestampToString(deadline)
	msg.ControlAction = tlv.Splice(0x20, dl)
	return msg.splice()
}

// 手动清洗开关
func (msg *whControlMessage) Clean(status int) string {
	msg.ControlAction = tlv.Splice(0x1f, strconv.Itoa(status))
	return msg.splice()
}

// 清零报文
func (msg *whControlMessage) Clear(status int8) string {
	if status&0x01 == 0x01 {
		msg.ControlAction = tlv.Splice(0x39, strconv.Itoa(0))
	} else if status&0x02 == 0x02 {
		msg.ControlAction = tlv.Splice(0x38, strconv.Itoa(0))
	} else if status&0x04 == 0x04 {
		msg.ControlAction = tlv.Splice(0x37, strconv.Itoa(0))
	} else if status&0x08 == 0x08 {
		msg.ControlAction = tlv.Splice(0x36, strconv.Itoa(0))
	} else if status&0x10 == 0x10 {
		msg.ControlAction = tlv.Splice(0x35, strconv.Itoa(0))
	}

	return msg.splice()
}