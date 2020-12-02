package tlv

import (
	"fmt"
)

const (
	HomeConsoleVersion = "Homeconsole04.00"
	packageName        = "tlv"
)

type TLV struct {
	Tag    int
	Length int
	Value  string
}

func (tlv *TLV) String() string {
	return fmt.Sprintf("%04X%04X%s", tlv.Tag, tlv.Length, tlv.Value)
}

// TLV长度，包含头部
func (tlv *TLV) Size() int {
	return tlv.Length + 8
}

// 拼接Homeconsole头部
func SpliceHead() string {
	s := HomeConsoleVersion + fmt.Sprintf("%08x", 1)
	return s
}

// 拼接TLV
// tag: 信元编码
// val: 数据
// 返回：编码后的字符串
func Splice(tag int, val string) string {
	return fmt.Sprintf("%04X%04X%s", tag, len(val), val)
}
