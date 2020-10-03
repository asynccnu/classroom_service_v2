package model

// 数据库表的结构
type Classrooms struct {
	Week                int
	Weekday             int
	Building            string `gorm:"type:char(1)"`
	AvailableClassrooms []byte `gorm:"type:json"`
}

// 参照上一版接口文档写的
type Rooms struct {
	One    []string
	Two    []string
	Three  []string
	Four   []string
	Five   []string
	Six    []string
	Seven  []string
	Eight  []string
	Nine   []string
	Ten    []string
	Eleven []string
	Twelve []string
}

type UnavailableClassrooms struct {
	Week     [2]int
	Weekday  int
	Time     [2]int
	WeekType string // 单双周
	Place    string
}

func InsertClassroomsInDB(instance *Classrooms) error {
	result := DB.Self.Create(instance)
	return result.Error
}

func UpdateAvailableClassroomInDB(instance *Classrooms) {
	DB.Self.Model(instance).Where("week = ? AND weekday= ? AND building= ?", instance.Week, instance.Weekday, instance.Building).Update("available_classrooms", instance.AvailableClassrooms)
}

func GetClassroomsFromDB(week int, weekday int, building string) (*Classrooms, error) {
	classroom := Classrooms{}
	result := DB.Self.Where("week = ? AND weekday = ? AND building = ?", week, weekday, building).First(&classroom)

	return &classroom, result.Error
}
