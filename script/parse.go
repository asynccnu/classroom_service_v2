package script

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/asynccnu/classroom_service_v2/log"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
)

// 从选课手册中解析获得课程信息
func GetCourseInfoFromCourseFile(channel chan *CourseItem, file *xlsx.File) {
	weekdayMap := map[string]int{"一": 1, "二": 2, "三": 3, "四": 4, "五": 5, "六": 6, "日": 7}

	// 正则
	r, err := regexp.Compile("星期(.*)第(.*)节{(.*)}")
	if err != nil {
		log.Error("Regexp compile failed", zap.String("reason", err.Error()))
		return
	}

	// 选课手册（0-15列）中第 10，11，12 列为上课时间，第 13，14，15 列为上课地点
	// 时间的格式：
	//    一般格式：星期一第9-10节{1-15周} 或 星期一第9-10节{1-15周(单)}
	//    特殊格式：星期一第9-10节{4-6周,8周} ...
	for _, sheet := range file.Sheets {
		fmt.Println(sheet.Name)

		// 遍历课程数据
		for i, row := range sheet.Rows {
			// 遍历一条课程数据中的多个时间、地点
			for j := 10; j <= 14; j += 2 {
				date := row.Cells[j].String()
				place := strings.ToUpper(row.Cells[j+1].String())
				if date == "" || place == "" {
					continue
				}
				// 跳过非7号楼、8号楼、南湖楼的课程
				if place[:1] != "7" && place[:1] != "8" && place[:1] != "N" {
					continue
				}

				// 正则匹配
				matches := r.FindStringSubmatch(date)
				if len(matches) < 4 {
					log.Error("Regexp match failed", zap.Int("row num", i))
					continue
				}

				weeks, err := ExtractWeeks(matches[3])
				if err != nil {
					log.Error("ExtractWeeks error", zap.Int("row num", i))
					continue
				}

				timeStart, timeEnd, err := ExtractClassTime(matches[2])
				if err != nil {
					log.Error("ExtractClassTime error", zap.Int("row num", i))
					continue
				}

				channel <- &CourseItem{
					Weeks: weeks,
					Day:   weekdayMap[matches[1]],
					Time:  [2]int{timeStart, timeEnd},
					Place: place,
				}
			}
		}
	}
	fmt.Println("Parsing course file OK")
}

// 处理节次
func ExtractClassTime(timeStr string) (int, int, error) {
	var classStart, classEnd int

	// 节次：1-2，1
	// 两种情况分别处理
	multiDuring := strings.Contains(timeStr, "-")
	if multiDuring {
		_, err := fmt.Sscanf(timeStr, "%d-%d", &classStart, &classEnd)
		if err != nil {
			return 0, 0, err
		}
	} else {
		_, err := fmt.Sscanf(timeStr, "%d", &classStart)
		if err != nil {
			return 0, 0, err
		}
		classEnd = classStart
	}

	return classStart, classEnd, nil
}

// 处理周次
func ExtractWeeks(weeksString string) ([]int, error) {
	var weeks []int

	doubleWeek := strings.Contains(weeksString, "双")
	singleWeek := strings.Contains(weeksString, "单")

	// 情况：逗号分隔的多个区间
	weekBlocks := strings.Split(weeksString, ",")
	for _, block := range weekBlocks {
		curWeeks, err := processWeeks(block, doubleWeek, singleWeek)
		if err != nil {
			log.Error("processWeeks function error")
			return nil, err
		}
		weeks = append(weeks, curWeeks...)
	}

	return weeks, nil
}

// 加工处理每一段的周次
func processWeeks(weekStr string, doubleWeek, singleWeek bool) ([]int, error) {
	var weeks []int
	var weekStart, weekEnd int

	// 1-10周，1周
	// 两种情况分别处理
	multiWeek := strings.Contains(weekStr, "-")
	if multiWeek {
		_, err := fmt.Sscanf(weekStr, "%d-%d", &weekStart, &weekEnd)
		if err != nil {
			log.Error("Split multiWeek error")
			return nil, err
		}
	} else {
		_, err := fmt.Sscanf(weekStr, "%d", &weekStart)
		if err != nil {
			return nil, err
		}
		weekEnd = weekStart
	}

	// 该字符串是否存在单双周标识，存在则以当前标识为准
	curDouble := strings.Contains(weekStr, "双")
	curSingle := strings.Contains(weekStr, "单")
	if curSingle || curDouble {
		doubleWeek = curDouble
		singleWeek = curSingle
	}

	for i := weekStart; i <= weekEnd; i++ {
		if doubleWeek && i%2 != 0 || singleWeek && i%2 == 0 {
			continue
		}
		weeks = append(weeks, i)
	}
	return weeks, nil
}
