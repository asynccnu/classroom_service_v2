package model

var DBName = "classroom"

const (
	ClassroomCol = "classroom"
)

type ClassroomModel struct {
	Week     int         `bson:"week" json:"week"`         // 周次
	Day      int         `bson:"day" json:"day"`           // 星期 1-7
	Building string      `bson:"building" json:"building"` // 教学楼，7/8/南湖综合楼
	List     []*RoomItem `bson:"list" json:"list"`
}

type RoomItem struct {
	Time  int      `bson:"time" json:"time"` // 节次，如第1节
	Rooms []string `bson:"rooms" json:"rooms"`
}
