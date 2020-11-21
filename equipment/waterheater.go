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

// 热水器设置状态
type WaterHeaterSetting struct {
	SerialNumber      	string
	SetActivateTime		int64
	Activate        	int8
	Unlock            	int8
	DeadlineTime    	int64
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

// 获取热水器设置状态
func (context *waterHeaterContext) LoadSetting(serialNumber string) (exists bool, setting *WaterHeaterSetting) {
	context.snapshot.Open()
	defer context.snapshot.Close()

	setting = new(WaterHeaterSetting)

	if !context.snapshot.Exists(waterHeaterPrefix + "setting_" + serialNumber) {
		return false, setting
	}

	context.snapshot.Load(waterHeaterPrefix + "setting_" + serialNumber, setting)

	return true, setting
}

// 保存热水器设置状态
func (context *waterHeaterContext) SaveSetting(setting *WaterHeaterSetting) {
	context.snapshot.Open()
	defer context.snapshot.Close()

	context.snapshot.Save(waterHeaterPrefix + "setting_" + setting.SerialNumber, setting)
}
