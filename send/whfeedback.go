package send

import (
	"github.com/robertzml/Gorgons/tlv"
	"strconv"
)

// 设备状态反馈报文
// 0x11
type wHFeedbackMessage struct {
	SerialNumber    string
	MainboardNumber string
	FeedbackAction   string
}

// 设备状态反馈报文构造函数
func NewWHFeedbackMessage(serialNumber string, mainboardNumber string) *wHFeedbackMessage{
	return &wHFeedbackMessage{ SerialNumber: serialNumber, MainboardNumber:mainboardNumber  }
}

// 拼接设备状态反馈报文
func (msg *wHFeedbackMessage) splice() string {
	head := tlv.SpliceHead()

	sn := tlv.Splice(0x127, msg.SerialNumber)
	mn := tlv.Splice(0x12b, msg.MainboardNumber)

	body := tlv.Splice(0x0011, sn + mn + msg.FeedbackAction)

	return head + body
}


// 快速响应
func (msg *wHFeedbackMessage) Fast(option int) string {
	msg.FeedbackAction = tlv.Splice(0x16, strconv.FormatInt(int64(option), 16))
	return msg.splice()
}

// 设备响应周期
func (msg *wHFeedbackMessage) Cycle(option int) string {
	msg.FeedbackAction = tlv.Splice(0x17, strconv.FormatInt(int64(option), 16))
	return msg.splice()
}

// 立即上报
func (msg *wHFeedbackMessage) Reply() string {
	msg.FeedbackAction = tlv.Splice(0x19, "1")
	return msg.splice()
}