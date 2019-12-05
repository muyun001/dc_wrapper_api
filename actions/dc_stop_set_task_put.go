package actions

import (
	"dc-wrapper-api/services/request_out_download_center"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DcStopSetTaskPut(c *gin.Context) {
	request_out_download_center.IS_SETTING_TASK = false

	c.JSON(http.StatusOK, gin.H{
		"msg": "停止发送任务",
	})
}
