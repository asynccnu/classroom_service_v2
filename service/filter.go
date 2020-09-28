package service

import (
	"github.com/asynccnu/classroom_service_v2/log"
	"github.com/asynccnu/classroom_service_v2/model"
	"runtime"

	"sync"

	"go.uber.org/zap"
)

var wg sync.WaitGroup

func InsertAllAndFilter() {

	// 先插入所有教室
	InsertAllClassrooms()
	wg.Add(8)
	channel := make(chan *model.UnavailableClassrooms, 8)

	for i := 0; i < 8; i++ {
		go RemoveUnavailableClassroomsInDB(channel)
	}
	go GetUnavailableClassrooms(channel)

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

			availableClassrooms := UnMarshalData(&classroomsInDB.AvailableClassrooms)
			for _, time := range unavailableClassrooms.Time {
				switch time {
				case 1:
					availableClassrooms.One = RemoveUnavailableClassroomFromList(&availableClassrooms.One, unavailableClassrooms.Place)
				case 2:
					availableClassrooms.Two = RemoveUnavailableClassroomFromList(&availableClassrooms.Two, unavailableClassrooms.Place)
				case 3:
					availableClassrooms.Three = RemoveUnavailableClassroomFromList(&availableClassrooms.Three, unavailableClassrooms.Place)
				case 4:
					availableClassrooms.Four = RemoveUnavailableClassroomFromList(&availableClassrooms.Four, unavailableClassrooms.Place)
				case 5:
					availableClassrooms.Five = RemoveUnavailableClassroomFromList(&availableClassrooms.Five, unavailableClassrooms.Place)
				case 6:
					availableClassrooms.Six = RemoveUnavailableClassroomFromList(&availableClassrooms.Six, unavailableClassrooms.Place)
				case 7:
					availableClassrooms.Seven = RemoveUnavailableClassroomFromList(&availableClassrooms.Seven, unavailableClassrooms.Place)
				case 8:
					availableClassrooms.Eight = RemoveUnavailableClassroomFromList(&availableClassrooms.Eight, unavailableClassrooms.Place)
				case 9:
					availableClassrooms.Nine = RemoveUnavailableClassroomFromList(&availableClassrooms.Nine, unavailableClassrooms.Place)
				case 10:
					availableClassrooms.Ten = RemoveUnavailableClassroomFromList(&availableClassrooms.Ten, unavailableClassrooms.Place)
				case 11:
					availableClassrooms.Eleven = RemoveUnavailableClassroomFromList(&availableClassrooms.Eleven, unavailableClassrooms.Place)
				case 12:
					availableClassrooms.Twelve = RemoveUnavailableClassroomFromList(&availableClassrooms.Twelve, unavailableClassrooms.Place)
				default:
					log.Error("Can't find right time")
				}

			}

			classroomsInDB.AvailableClassrooms = *MarshalData(availableClassrooms)

			model.UpdateAvailableClassroomInDB(classroomsInDB)
			availableClassrooms=nil
			classroomsInDB=nil
		}
		unavailableClassrooms=nil
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
