package actions

import (
	"dc-wrapper-api/databases"
	"dc-wrapper-api/structs/models"
	"dc-wrapper-api/structs/models/dc_request_logics"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// 录入请求任务
func RequestPut(c *gin.Context) {
	uniqueKey := c.Param("unique-key")
	dcRequest := models.DcRequest{}
	err := c.BindJSON(&dcRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "传入的Request格式错误",
		})
		return
	}
	dcRequest.UniqueKey = uniqueKey

	mgoSession := databases.MgoSession.Copy()
	defer mgoSession.Close()
	co := mgoSession.DB("fxt").C("fxt_requests")

	uniqueKeyIndex := mgo.Index{
		Key:        []string{"unique_key"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = co.EnsureIndex(uniqueKeyIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建unique_key索引失败",
		})
		return
	}

	autoExpireIndex := mgo.Index{
		Key:         []string{"updated_at"},
		ExpireAfter: time.Hour * 1,
	}
	err = co.EnsureIndex(autoExpireIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建updated_at索引失败",
		})
		return
	}

	statusIndex := mgo.Index{
		Key: []string{"status"},
	}
	err = co.EnsureIndex(statusIndex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "创建status索引失败",
		})
		return
	}

	existsDcReqeust := models.DcRequest{}
	err = co.Find(bson.M{"unique_key": dcRequest.UniqueKey}).One(&existsDcReqeust)
	if err != nil {
		dcRequest.CreatedAt = time.Now()
		dcRequest.UpdatedAt = time.Now()
		err = co.Insert(&dcRequest)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "插入记录失败",
			})
			return
		}
	} else {
		_ = co.Update(
			bson.M{
				"unique_key": dcRequest.UniqueKey,
				"status":     dc_request_logics.STATUS_已读取,
			},
			bson.M{
				"$set": bson.M{
					"status":     dc_request_logics.STATUS_未查询,
					"updated_at": time.Now(),
				},
			},
		)
	}

	c.JSON(http.StatusOK, dcRequest)
}
