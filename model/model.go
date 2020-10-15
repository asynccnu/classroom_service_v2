package model

var DBName = "classroom"

const (
	ClassroomCol = "classroom"
)

type ClassroomModel struct {
	Week     int         `bson:"week"`     // 周次
	Day      int         `bson:"day"`      // 星期 1-7
	Building string      `bson:"building"` // 教学楼，7/8/南湖综合楼
	List     []*RoomItem `bson:"list"`
}

type RoomItem struct {
	Time  int      `bson:"time"` // 节次，如第1节
	Rooms []string `bson:"rooms"`
}
