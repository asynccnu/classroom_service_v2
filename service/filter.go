package service

import (
	"runtime"

	"github.com/asynccnu/classroom_service_v2/log"
	"github.com/asynccnu/classroom_service_v2/model"

	"sync"

	"go.uber.org/zap"
)

var wg sync.WaitGroup

func InsertAllAndFilter(filePath string) {

	// 先插入所有教室
	InsertAllClassrooms()
	wg.Add(8)
	channel := make(chan *model.UnavailableClassrooms, 8)

	for i := 0; i < 8; i++ {
		go RemoveUnavailableClassroomsInDB(channel)
	}
	GetUnavailableClassrooms(channel, filePath)

	wg.Wait()
}

// 从数据库中移除被使用的教室
func RemoveUnavailableClassroomsInDB(channel chan *model.UnavailableClassrooms) {
	defer wg.Done()
	for {
		unavailableClassrooms, ok := <-channel
		if !ok {
			return
		}
		weekList := make([]int, 0)
		if unavailableClassrooms.WeekType == "" {
			for i := unavailableClassrooms.Week[0]; i <= unavailableClassrooms.Week[1]; i++ {
				weekList = append(weekList, i)
			}
		} else {
			for i := unavailableClassrooms.Week[0]; i <= unavailableClassrooms.Week[1]; i += 2 {
				weekList = append(weekList, i)
			}
		}

		for _, week := range weekList {
			classroomsInDB, err := model.GetClassroomsFromDB(week, unavailableClassrooms.Weekday, unavailableClassrooms.Place[:1])
			if err != nil {
				log.Error("GetClassroomsFromDB failed",
					zap.String("reason", err.Error()),
				)
				continue
			}

			// 这个取址符要加,我以为不加也行的
			availableClassrooms := &classroomsInDB.AvailableClassrooms

			classroomMap := map[int]*[]string{1: &availableClassrooms.One, 2: &availableClassrooms.Two, 3: &availableClassrooms.Three,
				4: &availableClassrooms.Four, 5: &availableClassrooms.Five, 6: &availableClassrooms.Six, 7: &availableClassrooms.Seven,
				8: &availableClassrooms.Eight, 9: &availableClassrooms.Nine, 10: &availableClassrooms.Ten,
				11: &availableClassrooms.Eleven, 12: &availableClassrooms.Twelve}

			for _, time := range unavailableClassrooms.Time {
				classroomWithTime, exist := classroomMap[time]
				if !exist {
					log.Error("wrong classroom Time",
					)
					continue
				}

				*classroomWithTime = RemoveUnavailableClassroomFromList(classroomWithTime, unavailableClassrooms.Place)
			}

			err = model.UpdateAvailableClassroomInDB(classroomsInDB)
			if err != nil {
				log.Error("Update available classroom in db failed",
					zap.String("reason", err.Error()),
				)
			}

		}

		runtime.GC()
	}

}

// 从教室列表中剔除被使用的教室
func RemoveUnavailableClassroomFromList(availableRooms *[]string, unavailableClassroom string) []string {
	newAvailableRooms := make([]string, 0)
	for _, room := range *availableRooms {
		if room == unavailableClassroom {
			continue
		}
		newAvailableRooms = append(newAvailableRooms, room)
	}

	return newAvailableRooms
}
