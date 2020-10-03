package classroom

import (
	"strconv"

	"github.com/asynccnu/classroom_service_v2/handler"
	"github.com/asynccnu/classroom_service_v2/model"
	"github.com/asynccnu/classroom_service_v2/pkg/errno"
	"github.com/asynccnu/classroom_service_v2/service"

	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	// 周数(1-20) 星期(1-7) 楼栋"7""8""N"
	week := c.DefaultQuery("weekno", "")
	weekday := c.DefaultQuery("week", "")
	building := c.DefaultQuery("building", "")
	if week == "" || weekday == "" || building == "" {
		handler.SendBadRequest(c, errno.ErrQuery, nil, "weekno, week and building are required.")
		return
	}

	weeknoString, _ := strconv.Atoi(week)
	weekString, _ := strconv.Atoi(weekday)
	classroom, err := model.GetClassroomsFromDB(weeknoString, weekString, building)
	if err != nil {
		handler.SendError(c, errno.ErrGetAvailableClassrooms, nil, err.Error())
		return
	}

	availableClassrooms := service.UnMarshalData(&classroom.AvailableClassrooms)

	handler.SendResponse(c, nil, availableClassrooms)
}
