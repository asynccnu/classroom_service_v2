package classroom

import (
	"github.com/asynccnu/classroom_service_v2/handler"
	"github.com/asynccnu/classroom_service_v2/service"
	"github.com/gin-gonic/gin"
)

// 重新更新数据库数据,一学期调用一次的那种
func Refresh(c *gin.Context){
	service.InsertAllAndFilter()
	handler.SendResponse(c,nil,"success")
}
