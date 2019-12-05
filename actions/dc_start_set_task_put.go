package actions

import (
	"dc-wrapper-api/services/request_out_download_center"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DcStartSetTaskPut(c *gin.Context) {
	request_out_download_center.IS_SETTING_TASK = true

	c.JSON(http.StatusOK, gin.H{
		"msg": "开始发送任务",
	})
}
