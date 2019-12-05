package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// 重置请求任务
func RequestResetPut(c *gin.Context) {
	uniqueKey := c.Param("unique-key")
	if uniqueKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请指定UniqueKey",
		})
		return
	}

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	err := co.Update(
		bson.M{"unique_key": uniqueKey},
		bson.M{
			"$set": bson.M{
				"status":     dc_request_logics.STATUS_未查询,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "已重置",
	})
}
