package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// 读取任务返回
func ResponseGet(c *gin.Context) {
	uniqueKey := c.Param("unique-key")

	redisSession := databases.RedisSession

	result, err := redisSession.Get(uniqueKey).Result()
	if err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "暂无结果",
		})
		return
	}

	dcResponse := models.DcResponse{}

	err = json.Unmarshal([]byte(result), &dcResponse)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "结果解析失败",
		})
		return
	}

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")
	err = co.Update(bson.M{"unique_key": uniqueKey}, bson.M{"$set": bson.M{"status": dc_request_logics.STATUS_已读取}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "更新状态失败",
		})
		return
	}

	redisSession.Del(uniqueKey)

	c.JSON(http.StatusOK, dcResponse)
}
