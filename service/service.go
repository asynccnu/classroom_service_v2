package service

import (
	"encoding/json"
	"github.com/asynccnu/classroom_service_v2/log"
	"github.com/asynccnu/classroom_service_v2/model"
	"github.com/tealeg/xlsx/v3"
	"go.uber.org/zap"

	"regexp"
	"strconv"
)

// 选课手册地址
var filePath = "/home/wency/data/2020—2021学年第一学期选课手册.xlsx"

func InsertAllClassrooms() {
	// 7,8号楼和南湖所有教室
	allClassroomsIn7 := []string{"7101", "7102", "7103", "7104", "7105", "7106", "7107", "7108", "7109",
		"7201", "7202", "7203", "7204", "7205", "7206", "7207", "7208", "7209", "7211",
		"7301", "7302", "7303", "7304", "7305", "7306", "7307", "7308", "7309", "7311",
		"7401", "7402", "7403", "7404", "7405", "7406", "7407", "7408", "7409", "7410", "7411",
		"7501", "7503", "7505",
	}

	allClassroomsIn8 := []string{
		"8101", "8102", "8103", "8104", "8105", "8106", "8107", "8108", "8109",
		"8110", "8111", "8112", "8201", "8202", "8203", "8204", "8205", "8206",
		"8207", "8208", "8209", "8210", "8211", "8212", "8213", "8214", "8215",
		"8216", "8301", "8302", "8303", "8304", "8305", "8306", "8307", "8308",
		"8309", "8310", "8311", "8312", "8313", "8314", "8315", "8316", "8401",
		"8402", "8403", "8404", "8405", "8406", "8407", "8408", "8409", "8410",
		"8411", "8412", "8413", "8414", "8415", "8416", "8501", "8502", "8503",
		"8504", "8505", "8506", "8507", "8508", "8509", "8510", "8511", "8512",
		"8513", "8514", "8515", "8516", "8716", "8717",
	}

	//NH都表示南湖
	allClassroomsInNH := []string{
		"N101", "N102", "N103", "N104", "N105", "N106", "N107", "N108", "N109", "N110", "N111", "N112", "N115", "N117", "N119",
		"N201", "N202", "N203", "N204", "N205", "N206", "N207", "N208", "N209", "N210", "N211", "N212", "N213", "N214", "N215", "N216", "N217", "N219", "N221", "N223",
		"N301", "N302", "N303", "N304", "N305", "N306", "N307", "N308", "N309", "N310", "N311", "N312", "N313",
		"N314", "N315", "N316", "N317", "N318", "N319", "N320", "N321", "N323", "N325", "N327",
	}

	building7 := model.Rooms{
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
		allClassroomsIn7,
	}

	building8 := model.Rooms{
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
		allClassroomsIn8,
	}

	buildingNH := model.Rooms{
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
		allClassroomsInNH,
	}

	// 先往数据库中插入所有教室  有的课周六都要上 丧心病狂
	for week := 1; week <= 21; week++ {
		for weekday := 1; weekday <= 7; weekday++ {
			instance7 := model.Classrooms{Week: week, Weekday: weekday, Building: "7", AvailableClassrooms: *MarshalData(&building7)}
			instance8 := model.Classrooms{Week: week, Weekday: weekday, Building: "8", AvailableClassrooms: *MarshalData(&building8)}
			instanceNH := model.Classrooms{Week: week, Weekday: weekday, Building: "N", AvailableClassrooms: *MarshalData(&buildingNH)}

			InsertData(&instance7)
			InsertData(&instance8)
			InsertData(&instanceNH)
		}
	}

}

// 从选课手册中解析被使用的教室
func GetUnavailableClassrooms(channel chan *model.UnavailableClassrooms) {
	weekdayMap := map[string]int{"一": 1, "二": 2, "三": 3, "四": 4, "五": 5,"六":6,"七":7}
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Fatal("Open xlsxFlie failed",
			zap.String("reason", err.Error()),
		)
	}

	// 选课手册中第10,11 \12,13\ 14,15 列为上课时间地点
	// 时间的格式为星期一第9-10节{1-15周(单)} 或 星期一第9-10节{1-15周} 或 星期一第9-10节{1-15周}
	for _, sheet := range file.Sheets {
		max:=sheet.MaxRow
		// 不要把下面max替换成sheet.MaxRow, 会死机,我吐了
		for i := 0; i <= max; i++ {
			for j := 10; j <= 14; j += 2 {
				cellDate, _ := sheet.Cell(i, j)
				cellPlace, _ := sheet.Cell(i, j+1)
				if cellDate.Value != "" {
					date := cellDate.String()
					place := cellPlace.String()
					if place==""{
						continue
					}
					if place[:1] == "7" || place[:1] == "8" || place[:1] == "N" {
						// 写正则匹配  匹配结构 [[第5-6节{6-20周(双)} 5 6 6 20 (双)]]
						r, _ := regexp.Compile("星期(.*)第(.*)-(.*)节{(.*)-(.*)周(.*)}")
						result := r.FindAllStringSubmatch(date, -1)

						// 当格式有变 没匹配到正确的数据时跳过
						if len(result) < 1 || len(result[0]) < 6 {
							continue
						}

						weekStart, _ := strconv.Atoi(result[0][4])
						weekEnd, _ := strconv.Atoi(result[0][5])
						timeStart, _ := strconv.Atoi(result[0][2])
						timeEnd, _ := strconv.Atoi(result[0][3])
						weekType := ""
						if len(result[0]) > 6 {
							weekType = result[0][6]
						}
						unavailableClassrooms := model.UnavailableClassrooms{
							Week:     [2]int{weekStart, weekEnd},
							Weekday:  weekdayMap[result[0][1]],
							Time:     [2]int{timeStart, timeEnd},
							WeekType: weekType,
							Place:    place,
						}
						channel <- &unavailableClassrooms
					}
				}

			}

		}
		sheet.Close()
	}
	close(channel)
	return
}

// 转化为JSON格式
func MarshalData(rooms *model.Rooms) *[]byte {
	bytes, err := json.Marshal(*rooms)
	if err != nil {
		log.Error("Marshal data failed",
			zap.String("reason", err.Error()),
		)
	}
	return &bytes
}

// 解析JSON格式
func UnMarshalData(rooms *[]byte) *model.Rooms {
	classroom := model.Rooms{}
	err := json.Unmarshal(*rooms, &classroom)
	if err != nil {
		log.Error("UnMarshal data failed",
			zap.String("reason", err.Error()),
		)
	}
	return &classroom
}

// 往数据库中插入数据
func InsertData(classrooms *model.Classrooms) {
	err := model.InsertClassroomsInDB(classrooms)
	if err != nil {
		log.Error("Insert data Faild",
			zap.String("reason", err.Error()),
		)
	}
}
