package script

import (
	"fmt"

	"github.com/asynccnu/classroom_service_v2/log"
	"github.com/tealeg/xlsx"

	"go.uber.org/zap"
)

// 解析并导入空闲教室数据
func SyncImportClassroomData(filePath string) {
	// 打开 Excel 文件
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		log.Fatal("Open xlsxFlie failed", zap.String("reason", err.Error()))
		return
	}

	// 先插入所有教室
	InsertAllClassrooms()

	channel := make(chan *CourseItem, 10)

	// 解析获取获取课程信息
	go func() {
		defer close(channel)
		GetCourseInfoFromCourseFile(channel, file)
	}()

	// 根据课程，从数据中移除被占用的教室
	for item := range channel {
		RemoveBusyRoomsByCourseInfo(item)
	}
	fmt.Println("Remove busy classrooms OK")

	// 导入空闲教室数据至数据库
	ImportDataToDB()

	fmt.Println("Import data into DB OK")
}
