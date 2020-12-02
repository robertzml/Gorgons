package send

import (
	"github.com/robertzml/Gorgons/tlv"
)

// 设备特殊控制报文
// 0x10 或 0x11
type wHSpecialMessage struct {
	SerialNumber    string
	MainboardNumber string
	SpecialAction   string
}

// 设备特殊控制报文构造函数
func NewWHSpecialMessage(serialNumber string, mainboardNumber string) *wHSpecialMessage {
	return &wHSpecialMessage{SerialNumber: serialNumber, MainboardNumber: mainboardNumber}
}

// 软件功能报文
func (msg *wHSpecialMessage) SoftFunction(option string) string {
	msg.SpecialAction = tlv.Splice(0x1d, option)

	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	sa := tlv.Splice(0x012, msg.SpecialAction)

	body := tlv.Splice(0x0010, sn+mn+sa)

	return head + body
}

// 热水器主控板特殊参数报文
func (msg *wHSpecialMessage) Special(option string) string {
	msg.SpecialAction = tlv.Splice(0x22, option)

	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	sa := tlv.Splice(0x012, msg.SpecialAction)

	body := tlv.Splice(0x0010, sn+mn+sa)

	return head + body
}

// 热水器手动控制拼写报文
func (msg *wHSpecialMessage) Manual(option string) string {
	msg.SpecialAction = option

	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)
	sa := tlv.Splice(0x012, msg.SpecialAction)

	body := tlv.Splice(0x0010, sn+mn+sa)

	return head + body
}

// 设备重复
// D8 设备序列号重复
// D7 主板序列号重复
func (msg *wHSpecialMessage) Duplicate(option string) string {
	msg.SpecialAction = tlv.Splice(0x13, option)

	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)

	body := tlv.Splice(0x0011, sn+mn+msg.SpecialAction)

	return head + body
}
