package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// 批量检查指定uniqueKey列表，返回其中已查完的uniqueKey列表
func ResponsesCheckPost(c *gin.Context) {
	var uniqueKeys []string
	err := c.BindJSON(&uniqueKeys)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "请传入JSON格式的uniqueKey数组",
		})
	}

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	// updated_at超过10min,重置状态为0
	_ = co.Update(
		bson.M{
			"status":     dc_request_logics.STATUS_查询中,
			"updated_at": bson.M{"$lt": time.Now().Add(-time.Minute * 10)},
			"unique_key": bson.M{"$in": uniqueKeys},
		},
		bson.M{
			"$set": bson.M{
				"status":     dc_request_logics.STATUS_未查询,
				"updated_at": time.Now(),
			},
		},
	)

	// 数据已读,重置状态为0
	_ = co.Update(
		bson.M{
			"status":     dc_request_logics.STATUS_已读取,
			"unique_key": bson.M{"$in": uniqueKeys},
		},
		bson.M{
			"$set": bson.M{
				"status":     dc_request_logics.STATUS_未查询,
				"updated_at": time.Now(),
			},
		},
	)

	finishedFilter := bson.M{"status": dc_request_logics.STATUS_已查询}
	if len(uniqueKeys) > 0 {
		finishedFilter["unique_key"] = bson.M{"$in": uniqueKeys}
	}

	finishedUniqueKeys := make([]string, 0)
	dcRequest := models.DcRequest{}
	iter := co.Find(finishedFilter).Iter()
	for iter.Next(&dcRequest) {
		redisSession := databases.RedisSession
		keys, err := redisSession.Keys(dcRequest.UniqueKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Redis服务故障",
			})
			return
		}
		hasResponse := len(keys) != 0
		if hasResponse {
			finishedUniqueKeys = append(finishedUniqueKeys, dcRequest.UniqueKey)
		} else {
			dcRequest.Status = dc_request_logics.STATUS_未查询
			dcRequest.CreatedAt = time.Now()
			err = co.Update(bson.M{"unique_key": dcRequest.UniqueKey}, dcRequest)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "更新Request状态失败",
				})
				return
			}
		}
	}

	c.JSON(http.StatusOK, finishedUniqueKeys)
}
