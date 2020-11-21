package send

// 下发报文协议处理模块，主要实现协议转换，拼接成homeconsole
// 1. 热水器控制报文
// 2. 热水器状态反馈报文

const (
	packageName = "send"
)

// 报文消息接口
// 所有下发的报文均实现该接口
type Message interface {

	splice() string
}