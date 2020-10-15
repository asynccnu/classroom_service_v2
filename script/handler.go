package script

import (
	"github.com/asynccnu/classroom_service_v2/log"
	"github.com/asynccnu/classroom_service_v2/model"

	"go.uber.org/zap"
)

const (
	weekMin = 4  // 最小周数
	weekMax = 21 // 最大周数
	dayMin  = 1
	dayMax  = 7
)

var instances []*model.ClassroomModel

// 解析到的课程信息，借此移除被占用的教室
type CourseItem struct {
	Weeks []int
	Day   int    // 星期
	Time  [2]int // 节次，分别为开始和结束
	Place string
}

// 暂存所有教室数据
func InsertAllClassrooms() {
	// 7，8号楼和南湖综合楼的所有教室
	var allClassrooms = [3][]string{
		{
			"7101", "7102", "7103", "7104", "7105", "7106", "7107", "7108", "7109",
			"7201", "7202", "7203", "7204", "7205", "7206", "7207", "7208", "7209", "7211",
			"7301", "7302", "7303", "7304", "7305", "7306", "7307", "7308", "7309", "7311",
			"7401", "7402", "7403", "7404", "7405", "7406", "7407", "7408", "7409", "7410", "7411",
			"7501", "7503", "7505",
		}, {
			"8101", "8102", "8103", "8104", "8105", "8106", "8107", "8108", "8109",
			"8110", "8111", "8112", "8201", "8202", "8203", "8204", "8205", "8206",
			"8207", "8208", "8209", "8210", "8211", "8212", "8213", "8214", "8215",
			"8216", "8301", "8302", "8303", "8304", "8305", "8306", "8307", "8308",
			"8309", "8310", "8311", "8312", "8313", "8314", "8315", "8316", "8401",
			"8402", "8403", "8404", "8405", "8406", "8407", "8408", "8409", "8410",
			"8411", "8412", "8413", "8414", "8415", "8416", "8501", "8502", "8503",
			"8504", "8505", "8506", "8507", "8508", "8509", "8510", "8511", "8512",
			"8513", "8514", "8515", "8516",
		}, {
			"N101", "N102", "N103", "N104", "N105", "N106", "N107", "N108", "N109", "N110", "N111", "N112", "N115", "N117", "N119",
			"N201", "N202", "N203", "N204", "N205", "N206", "N207", "N208", "N209", "N210", "N211", "N212", "N213",
			"N214", "N215", "N216", "N217", "N219", "N221", "N223",
			"N301", "N302", "N303", "N304", "N305", "N306", "N307", "N308", "N309", "N310", "N311", "N312", "N313",
			"N314", "N315", "N316", "N317", "N318", "N319", "N320", "N321", "N323", "N325", "N327",
		},
	}

	buildings := [3]string{"7", "8", "N"}

	// 先插入所有教室
	// 遍历周次
	for week := weekMin; week <= weekMax; week++ {
		// 遍历星期
		for day := dayMin; day <= dayMax; day++ {
			// 楼栋
			for i := 0; i < 3; i++ {
				var roomList = make([]*model.RoomItem, 0)

				// 全部节次：1-12
				for num := 1; num <= 12; num++ {
					rooms := make([]string, len(allClassrooms[i]))
					copy(rooms, allClassrooms[i][:])
					roomList = append(roomList, &model.RoomItem{
						Time:  num,
						Rooms: rooms,
					})
				}

				instances = append(instances, &model.ClassroomModel{
					Week:     week,
					Day:      day,
					Building: buildings[i],
					List:     roomList,
				})
			}
		}
	}
}

// 移除被占用的教室
func RemoveBusyRoomsByCourseInfo(course *CourseItem) {

	// 需要移除的教室
	place := course.Place
	building := place[:1]

	weekMap := make(map[int]bool, len(course.Weeks))
	for _, week := range course.Weeks {
		weekMap[week] = true
	}

	// 遍历空闲教室暂存记录
	for _, instance := range instances {
		// 楼栋、星期是否一致
		if instance.Building != building || instance.Day != course.Day {
			continue
		}

		// 是否是该周的
		if _, ok := weekMap[instance.Week]; !ok {
			continue
		}

		// 遍历该天全部节次
		for _, roomItem := range instance.List {
			// 是否在上课节次区间内
			if roomItem.Time < course.Time[0] || roomItem.Time > course.Time[1] {
				continue
			}

			// 遍历该节课的教室
			for k, room := range roomItem.Rooms {
				// 移除被占用的教室
				// 目前只有一个需要移除
				if room == place {
					// 切片底层共享一个数组，需要新建一个数组，避免对其它造成连带修改
					// s := make([]string, len(roomItem.Rooms)-1)
					// copy(s, roomItem.Rooms[:k])
					// copy(s[k:], roomItem.Rooms[k+1:])
					// instances[i].List[j].Rooms = s
					roomItem.Rooms = append(roomItem.Rooms[:k], roomItem.Rooms[k+1:]...)
					break
				}
			}
		}
	}
}

// 导入数据至数据库
func ImportDataToDB() {
	// 批量插入数据
	if err := model.CreateMultipleClassroomDocs(instances); err != nil {
		log.Error("Inserting multiple data failed", zap.String("reason", err.Error()))
	}
}
