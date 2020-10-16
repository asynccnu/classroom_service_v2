package classroom

import (
	"strconv"

	"github.com/asynccnu/classroom_service_v2/handler"
	"github.com/asynccnu/classroom_service_v2/model"
	"github.com/asynccnu/classroom_service_v2/pkg/errno"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	// 周数(1-20) 星期(1-7) 楼栋"7""8""N"
	weekStr := c.DefaultQuery("week", "")
	dayStr := c.DefaultQuery("day", "")
	building := c.DefaultQuery("building", "")
	if weekStr == "" || dayStr == "" || building == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "The week, day and building are required.")
		return
	}

	week, err := strconv.Atoi(weekStr)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "The 'week' is wrong.")
		return
	}

	day, err := strconv.Atoi(dayStr)
	if err != nil {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "The 'day' is wrong.")
		return
	}

	// 获取空闲教室
	classroom, err := model.GetClassroomDoc(week, day, building)
	if mongo.ErrNoDocuments == err {
		classroom = &model.ClassroomModel{
			Week:     week,
			Day:      day,
			Building: building,
			List:     make([]*model.RoomItem, 0),
		}
	} else if err != nil {
		handler.SendError(c, errno.ErrGetClassrooms, nil, err.Error())
		return
	}

	handler.SendResponse(c, nil, classroom)
}
