package equipment

import "github.com/robertzml/Gorgons/db"

const (
	// 热水器Redis前缀
	waterHeaterPrefix = "wh_"
)

// 热水器数据处理类
type waterHeaterContext struct {
	// 实时数据操作接口
	snapshot 	db.Snapshot
}

func NewWaterHeaterContext(snap db.Snapshot) *waterHeaterContext {
	context := new(waterHeaterContext)
	context.snapshot = snap

	return context
}

// 通过设备序列号获取主板序列号
func (context *waterHeaterContext) GetMainboardNumber(serialNumber string) (mainboardNumber string, exists bool) {
	context.snapshot.Open()
	defer context.snapshot.Close()

	mn := context.snapshot.LoadField(waterHeaterPrefix + serialNumber, "MainboardNumber")
	if len(mn) == 0 {
		return "",false
	} else {
		return mn,true
	}
}